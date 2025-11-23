// Package config: contains all the configurations run the api
package config

import (
	"encoding/json"
	"os"

	"github.com/sunimalherath/orderfoodonline/internal/app/utils"
	"github.com/sunimalherath/orderfoodonline/internal/core/entities"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port string
}

func Load() *Config {
	return &Config{
		ServerConfig{
			Port: utils.GetEnvVar("PORT", "8080"),
		},
	}
}

func LoadProducts() (map[string]entities.Product, error) {
	prodData, err := os.ReadFile("./internal/config/data/products.json")
	if err != nil {
		return nil, err
	}

	var prodCache map[string]entities.Product

	if err := json.Unmarshal(prodData, &prodCache); err != nil {
		return nil, err
	}

	return prodCache, nil
}
