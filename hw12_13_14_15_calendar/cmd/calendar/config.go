package main

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf `yaml:"log"`
	// TODO
}

type LoggerConf struct {
	Level string `yaml:"level"`
	// TODO
}

func NewConfig() Config {
	return Config{}
}

func LoadConfigFile(cfg interface{}, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, cfg)
}
