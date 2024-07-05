package main

import (
	"flag"
	"gopkg.in/yaml.v3"
	"hua-cloud.com/tools/image-sync/internal/config"
	"hua-cloud.com/tools/image-sync/internal/server"
	"os"
)

var configFile string

func main() {
	// 加载配置
	cfg := loadConfig()
	// 创建服务
	imageServer := server.NewImageServer(cfg)
	// 启动服务
	imageServer.Start()
}

func init() {
	flag.StringVar(&configFile, "f", "config.yaml", "config file")
}

func loadConfig() *config.ServerConfig {
	flag.Parse()
	file, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	cfg := config.ServerConfig{}
	if err = yaml.Unmarshal(file, &cfg); err != nil {
		panic(err)
	}
	return &cfg
}
