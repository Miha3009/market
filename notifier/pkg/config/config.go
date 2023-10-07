package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type EmailConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	From     string `yaml:"from"`
}

type KafkaConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Topic     string `yaml:"topic"`
	Partition int    `yaml:"partition"`
}

type Config struct {
	Server EmailConfig `yaml:"emailServer"`
	Kafka  KafkaConfig `yaml:"kafka"`
}

func ReadConfig(path string) (*Config, error) {
	cfg := &Config{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
