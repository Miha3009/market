package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port int `yaml:"port"`
}

type MongoDBConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
	Name string `yaml:"name"`
}

type InventoryConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type KafkaConfig struct {
	Port  int    `yaml:"port"`
	Host  string `yaml:"host"`
	Topic string `yaml:"topic"`
}

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  MongoDBConfig   `yaml:"mongo"`
	Inventory InventoryConfig `yaml:"inventory"`
	Kafka     KafkaConfig     `yamk:"kafka"`
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
