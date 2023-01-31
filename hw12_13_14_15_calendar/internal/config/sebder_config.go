package config

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

type SenderConfig struct {
	Logger   LoggerConf   `yaml:"log"`
	Consumer ConsumerConf `yaml:"consumer"`
}

type ConsumerConf struct {
	URI          string `yaml:"uri"`
	Exchange     string `yaml:"exchange"`
	ExchangeType string `yaml:"exchangeType"`
	Queue        string `yaml:"queue"`
	BindingKey   string `yaml:"bindingKey"`
	ConsumerTag  string `yaml:"consumerTag"`
	Lifetime     uint64 `yaml:"lifetime"`
}

func NewSenderConfig() SenderConfig {
	return SenderConfig{}
}

func LoadSenderConfig(filePath string) (SenderConfig, error) {
	cfg := NewSenderConfig()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
