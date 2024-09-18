package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/opencontainers/go-digest"
	"github.com/szpinc/image-sync/pkg/types"
)

var ErrUnauthorized = errors.New("Unauthorized")

type Client struct {
	config *types.ClientConfig
}

func NewClient(config *types.ClientConfig) *Client {
	return &Client{config: config}
}

func (c *Client) CheckBlobExists(repository string, digest digest.Digest) (bool, error) {

	result, err := c.doRequest(
		http.MethodGet,
		fmt.Sprintf("%s/api/v1/blob/exists?repository=%s&digest=%s", c.config.Server.Address, repository, digest),
		nil,
		"application/json",
	)

	if err != nil {
		return false, checkError(err)
	}

	exists, _ := result.Data.(bool)

	return exists, nil
}

func (c *Client) UploadBlob(repository string, digest digest.Digest, data io.Reader) error {

	url := fmt.Sprintf("%s/api/v1/blob/uploads?repository=%s&digest=%s", c.config.Server.Address, repository, digest)

	_, err := c.doRequest(http.MethodPut, url, data, "application/octet-stream")

	return checkError(err)
}

func (c *Client) PutManifest(repository string, tag string, manifest *schema2.DeserializedManifest) error {

	requestBody, _ := manifest.MarshalJSON()

	url := fmt.Sprintf("%s/api/v1/manifest/push?repository=%s&tag=%s", c.config.Server.Address, repository, tag)

	_, err := c.doRequest(http.MethodPost, url, bytes.NewBuffer(requestBody), "application/json")

	return checkError(err)
}

func (c *Client) Copy(srcRepository string, targetRepository string, tag string) error {

	// 创建registry client
	registryClient, err := newRegistryClient(c.config.Registry)

	if err != nil {
		return err
	}

	manifest, err := registryClient.ManifestV2(srcRepository, tag)

	if err != nil {
		return err
	}

	digests := []digest.Digest{}

	// 添加配置blob
	digests = append(digests, manifest.Config.Digest)

	for _, layer := range manifest.Layers {
		digests = append(digests, layer.Digest)
	}

	wg := sync.WaitGroup{}

	for _, blob := range digests {
		wg.Add(1)
		go func(blob digest.Digest) {
			defer wg.Done()
			// 上传layer,引入重试,5秒重试一次
			err := retry.Do(func() error {
				fmt.Printf("copy blob: %s\n", blob.String())
				return checkError(uploadBlob(srcRepository, targetRepository, blob, c, registryClient))
			},
				retry.Attempts(5),
				retry.Delay(time.Second*5),
			)
			if err != nil {
				panic(err)
			}
		}(blob)
	}
	wg.Wait()

	// 更新manifest
	if err = c.PutManifest(targetRepository, tag, manifest); err != nil {
		return err
	}

	return nil
}

func (c *Client) Deploy(request types.CmdRequest) error {
	url := fmt.Sprintf("%s/api/v1/deploy", c.config.Server.Address)
	reqBody, err := json.Marshal(request)

	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.config.Server.Username, c.config.Server.Password))))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	r, err := getRespBody(resp.Body)

	if err != nil {
		return err
	}

	if r.Code != http.StatusOK {
		return errors.New(r.Message)
	}
	return nil
}

func uploadBlob(srcRepo string, targetRepo string, digest digest.Digest, c *Client, srcRegistry *registry.Registry) error {
	// 检测blob hash是否存在
	exists, err := c.CheckBlobExists(targetRepo, digest)

	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	blob, err := srcRegistry.DownloadBlob(srcRepo, digest)

	if err != nil {
		return err
	}

	defer func(blob io.ReadCloser) {
		_ = blob.Close()
	}(blob)

	return c.UploadBlob(targetRepo, digest, blob)
}

func getRespBody(body io.ReadCloser) (*types.Resp, error) {

	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(body)

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	println(string(data))
	resp := types.Resp{}
	return &resp, json.Unmarshal(data, &resp)
}

func newRegistryClient(conf types.RegistryConfig) (*registry.Registry, error) {
	transport := registry.WrapTransport(http.DefaultTransport, conf.Url, conf.Username, conf.Password)

	url := strings.TrimSuffix(conf.Url, "/")

	reg := &registry.Registry{
		URL: url,
		Client: &http.Client{
			Transport: transport,
		},
		Logf: registry.Quiet,
	}

	if err := reg.Ping(); err != nil {
		return nil, err
	}

	return reg, nil
}

func (c *Client) doRequest(method string, url string, body io.Reader, contentType string) (*types.Resp, error) {

	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.config.Server.Username, c.config.Server.Password))))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	res := types.Resp{}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBody, &res)

	if err != nil {
		return nil, err
	}

	if res.Code != http.StatusOK {
		return &res, errors.New(res.Message)
	}

	return &res, nil
}

func checkError(err error) error {

	if err == nil {
		return nil
	}

	if err == ErrUnauthorized {
		fmt.Println("Unauthorized")
		os.Exit(1)
	}
	return err
}
