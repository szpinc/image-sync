package main

import (
	"hua-cloud.com/tools/image-sync/internal/config"
	"hua-cloud.com/tools/image-sync/internal/server"
)

func main() {

	imageServer := server.NewImageServer(&config.ServerConfig{
		Addr: ":23333",
		RegistryConfig: config.RegistryConfig{
			Url:      "https://registry.cn-hangzhou.aliyuncs.com",
			Username: "1316420259@qq.com",
			Password: "Szp11020652",
		},
		LogConfig: config.LogConfig{
			Level: "debug",
		},
	})

	imageServer.Start()
}
