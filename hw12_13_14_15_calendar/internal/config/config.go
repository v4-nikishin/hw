package config

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Logger LoggerConf   `yaml:"log"`
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

func NewConfig() Config {
	return Config{}
}

func Load(filePath string) (Config, error) {
	cfg := NewConfig()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, err
}
