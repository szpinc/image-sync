package server

import (
	"sync"
	"testing"

	"github.com/heroku/docker-registry-client/registry"
	"github.com/opencontainers/go-digest"
	"hua-cloud.com/tools/image-sync/internal/config"
)

func TestImageServer_GetManifest(t *testing.T) {

	repo, tag := "zhmz/monitor-service", "20240703170156-amd64"

	srcRegistry, err := registry.New("http://harbor.hy-zw.com", "admin", "Hyzs(%23@7&8")

	server := NewImageServer(&config.ServerConfig{
		Addr: ":23333",
		RegistryConfig: config.RegistryConfig{
			Url: "http://127.0.0.1:50000",
		},
		LogConfig: config.LogConfig{
			Level: "debug",
		},
	})

	manifest, err := srcRegistry.ManifestV2(repo, tag)

	if err != nil {
		t.Fatal(err)
	}

	// 上传配置清单
	if err := uploadBlob(repo, manifest.Config.Digest, server, srcRegistry); err != nil {
		t.Fatal(err)
	}

	t.Logf("manifest is: %v\n", manifest.Manifest)

	wg := sync.WaitGroup{}

	for _, layer := range manifest.Layers {

		wg.Add(1)

		func() {
			defer wg.Done()

			t.Logf("layer is: %v\n", layer.Digest)
			// 上传layer
			if err := uploadBlob(repo, layer.Digest, server, srcRegistry); err != nil {
				t.Fatal(err)
			}
		}()
	}

	wg.Wait()

	// 更新manifest
	if err = server.PushManifest(repo, tag, manifest); err != nil {
		t.Fatal(err)
	}
}

func uploadBlob(repo string, digest digest.Digest, serv *ImageServer, srcRegistry *registry.Registry) error {

	// 检测blob hash是否存在
	exists, err := serv.HasBlob(repo, digest)

	if err != nil {
		return err
	}

	if exists {
		serv.Log.Info("blog exists,ignored")
		return nil
	}

	blob, err := srcRegistry.DownloadBlob(repo, digest)

	if err != nil {
		return err
	}

	return serv.UploadBlob(digest, repo, blob)
}

func TestGetManifest(t *testing.T) {
	repo, tag := "park/tibmas2-webapi", "20240624114951"

	srcRegistry, _ := registry.New("http://harbor.hy-zw.com", "admin", "Hyzs(%23@7&8")

	manifest, err := srcRegistry.ManifestV2(repo, tag)

	if err != nil {
		t.Fatal(err)
	}
	for _, layer := range manifest.Layers {
		t.Log(layer.Digest)
	}
}
