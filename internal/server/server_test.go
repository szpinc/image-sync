package server

import (
	"github.com/heroku/docker-registry-client/registry"
	"hua-cloud.com/tools/image-sync/internal/config"
	"testing"
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

	if exists, err := server.HasBlob(repo, manifest.Config.Digest); err != nil {
		t.Fatal(err)
	} else if !exists {
		configBlob, err := srcRegistry.DownloadBlob(repo, manifest.Config.Digest)
		if err != nil {
			t.Fatal(err)
		}
		if err := server.UploadBlob(manifest.Config.Digest, repo, configBlob); err != nil {
			t.Fatal(err)
		}
	}

	t.Logf("manifest is: %v\n", manifest.Manifest)
	for _, layer := range manifest.Layers {
		t.Logf("layer is: %v\n", layer.Digest)
		exists, err := server.HasBlob(repo, layer.Digest)
		if err != nil {
			t.Fatal(err)
		}
		if exists {
			t.Logf("blog: %s exists: \n", layer.Digest.String())
			continue
		}
		blob, err := srcRegistry.DownloadBlob(repo, layer.Digest)

		if err != nil {
			t.Fatal(err)
		}
		err = server.UploadBlob(layer.Digest, repo, blob)
		if err != nil {
			t.Fatal(err)
		}
	}

	err = server.PushManifest(repo, tag, manifest)

	if err != nil {
		t.Fatal(err)
	}
}
