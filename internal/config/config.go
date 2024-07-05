package config

import "github.com/gin-gonic/gin"

type ServerConfig struct {
	Addr           string         `yaml:"addr"`
	RegistryConfig RegistryConfig `yaml:"registry"`
	LogConfig      LogConfig      `yaml:"log"`
	Accounts       gin.Accounts   `yaml:"accounts"`
}

type RegistryConfig struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}
