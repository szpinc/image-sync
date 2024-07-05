package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/opencontainers/go-digest"
	"hua-cloud.com/tools/image-sync/internal/types"
	"io"
	"net/http"
	"sync"
)

type Client struct {
	config *types.ClientConfig
}

func NewClient(config *types.ClientConfig) *Client {
	return &Client{config: config}
}

func (c *Client) CheckBlobExists(repository string, digest digest.Digest) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/blob/exists?repository=%s&digest=%s", c.config.Server.Address, repository, digest)

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.config.Server.Username, c.config.Server.Password))))

	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return false, errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return false, err
	}

	res := types.Resp{}

	_ = json.Unmarshal(body, &res)

	if res.Code != http.StatusOK {
		return true, errors.New(string(body))
	}

	exists, _ := res.Data.(bool)

	return exists, nil
}

func (c *Client) UploadBlob(repository string, digest digest.Digest, data io.Reader) error {

	url := fmt.Sprintf("%s/api/v1/blob/uploads?repository=%s&digest=%s", c.config.Server.Address, repository, digest)

	req, err := http.NewRequest(http.MethodPut, url, data)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.config.Server.Username, c.config.Server.Password))))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

func (c *Client) PutManifest(repository string, tag string, manifest *schema2.DeserializedManifest) error {

	requestBody, _ := manifest.MarshalJSON()

	url := fmt.Sprintf("%s/api/v1/manifest/push?repository=%s&tag=%s", c.config.Server.Address, repository, tag)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))

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
	return nil
}

func (c *Client) Copy(repository string, tag string) error {
	registryConf := c.config.Registry
	registryClient, err := registry.New(registryConf.Url, registryConf.Username, registryConf.Password)

	if err != nil {
		return err
	}
	registryClient.Logf = func(format string, args ...interface{}) {}

	manifest, err := registryClient.ManifestV2(repository, tag)

	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	go func() {
		wg.Add(1)
		defer wg.Done()
		// 上传配置清单
		if err := uploadBlob(repository, manifest.Config.Digest, c, registryClient); err != nil {
			panic(err)
		}
	}()

	for _, layer := range manifest.Layers {

		wg.Add(1)

		go func() {
			defer wg.Done()
			// 上传layer
			if err := uploadBlob(repository, layer.Digest, c, registryClient); err != nil {
				panic(err)
			}
		}()
	}

	wg.Wait()

	// 更新manifest
	if err = c.PutManifest(repository, tag, manifest); err != nil {
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

func uploadBlob(repo string, digest digest.Digest, c *Client, srcRegistry *registry.Registry) error {

	println("copy blob: ", digest.String())
	// 检测blob hash是否存在
	exists, err := c.CheckBlobExists(repo, digest)

	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	blob, err := srcRegistry.DownloadBlob(repo, digest)

	if err != nil {
		return err
	}

	defer func(blob io.ReadCloser) {
		_ = blob.Close()
	}(blob)

	return c.UploadBlob(repo, digest, blob)
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
