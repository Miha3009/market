//go:generate go run github.com/99designs/gqlgen generate

package graph

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v2"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type ProductConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type OrderConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ServerConfig struct {
	Port    int           `yaml:"port"`
	Product ProductConfig `yaml:"products"`
	Order   OrderConfig   `yaml:"orders"`
}

type Resolver struct {
	Client *resty.Client
	Cfg    *ServerConfig
}

func (r *Resolver) GetOrderService() string {
	return fmt.Sprintf("%s:%d", r.Cfg.Order.Host, r.Cfg.Order.Port)
}

func (r *Resolver) GetProductService() string {
	return fmt.Sprintf("%s:%d", r.Cfg.Product.Host, r.Cfg.Product.Port)
}

func ReadConfig(path string) (*ServerConfig, error) {
	cfg := &ServerConfig{}
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
