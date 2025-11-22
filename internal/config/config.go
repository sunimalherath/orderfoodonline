package config

import "github.com/sunimalherath/orderfoodonline/internal/app/utils"

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
