package types

type ClientConfig struct {
	Version  string         `yaml:"version"`
	Server   Server         `yaml:"server"`
	Registry RegistryConfig `yaml:"registry"`
}

type Server struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Auth     string `yaml:"auth"`
}

type RegistryConfig struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
