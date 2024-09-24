package client

import (
	"testing"

	"github.com/szpinc/image-sync/pkg/types"
)

func TestClient_CheckBlobExists(t *testing.T) {

	client := NewClient(&types.ClientConfig{
		Version: "v1",
		Server: types.Server{
			Address: "https://smartum.sz.gov.cn/park/image/sync",
		},
	})

	exists, err := client.CheckBlobExists("zhmz/monitor-service", "sha256:11d7d3a06ff9d48c44f0effc3ffa715e45c2e05701738f0769ffa870e8ea6fda")

	if err != nil {
		t.Fatal(err)
	}

	t.Log(exists)
}

func TestClient_Copy(t *testing.T) {

	//testcases := []struct {
	//	srcRepository    string
	//	targetRepository string
	//	tag              string
	//}{
	//	{
	//
	//	},
	//}

	client := NewClient(&types.ClientConfig{
		Version: "v1",
		Server: types.Server{
			Address:  "http://127.0.0.1:23333",
			Username: "admin",
			Password: "9YVtWmou1CxfJ3LXya8e",
		},
		Registry: types.RegistryConfig{Url: "http://harbor.hy-zw.com", Username: "admin", Password: "Hyzs(%23@7&8"},
	})

	err := client.Copy("zhmz/web-tibmas-oauth", "mzj_hyzs/web-tibmas-oauth", "20240920113435-arm64")

	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Deploy(t *testing.T) {
	client := NewClient(&types.ClientConfig{
		Version: "v1",
		Server: types.Server{
			Address: "https://smartum.sz.gov.cn/park/image/sync",
		},
		Registry: types.RegistryConfig{Url: "http://harbor.hy-zw.com", Username: "admin", Password: "Hyzs(%23@7&8"},
	})
	err := client.Deploy(types.CmdRequest{
		Host:              "10.226.22.6",
		Port:              36000,
		App:               "park-largescreen-web",
		Repository:        "default/park-largescreen-web",
		Tag:               "20240624114951",
		DockerComposeFile: "/data/docker-compose.yml",
		Deploy:            false,
	})
	if err != nil {
		t.Fatal(err)
	}
}
