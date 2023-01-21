package config

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

type SchedulerConfig struct {
	Logger    LoggerConf    `yaml:"log"`
	Server    ServerConf    `yaml:"server"`
	Scheduler SchedulerConf `yaml:"scheduler"`
	Publisher PublisherConf `yaml:"publisher"`
}

type SchedulerConf struct {
	CheckPeriod uint8 `yaml:"checkPeriod"`
}

type PublisherConf struct {
	URI          string `yaml:"uri"`
	Exchange     string `yaml:"exchange"`
	ExchangeType string `yaml:"exchangeType"`
	RoutingKey   string `yaml:"routingKey"`
	Reliable     bool   `yaml:"reliable"`
}

func NewSchedulerConfig() SchedulerConfig {
	return SchedulerConfig{}
}

func LoadSchedulerConfig(filePath string) (SchedulerConfig, error) {
	cfg := NewSchedulerConfig()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
