package config

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Logger LoggerConf `yaml:"log"`
	Server ServerConf `yaml:"server"`
	DB     DBConf     `yaml:"database"`
}

type ServerConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type DBConf struct {
	Type string `yaml:"type"`
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
