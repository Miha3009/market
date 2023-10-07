package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port int `yaml:"port"`
}

type PostgresConfig struct {
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type InventoryConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type RedisConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Postgres  PostgresConfig  `yaml:"postgres"`
	Inventory InventoryConfig `yaml:"inventory"`
	Redis     RedisConfig     `yaml:"redis"`
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
