// Package config: contains all the configurations run the api
package config

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/sunimalherath/orderfoodonline/internal/app/utils"
	"github.com/sunimalherath/orderfoodonline/internal/core/constants"
	"github.com/sunimalherath/orderfoodonline/internal/core/entities"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port string
}

func Load() *Config {
	loadEnvFile(constants.EnvFilePath)

	return &Config{
		ServerConfig{
			Port: utils.GetEnvVar(constants.PORT, "8080"),
		},
	}
}

func LoadProducts() (map[string]entities.Product, error) {
	prodData, err := os.ReadFile(constants.ProductsFilePath)
	if err != nil {
		return nil, err
	}

	var prodCache map[string]entities.Product

	if err := json.Unmarshal(prodData, &prodCache); err != nil {
		return nil, err
	}

	return prodCache, nil
}

func GetCouponFilePaths() []string {
	return []string{
		constants.CouponFilePath1,
		constants.CouponFilePath2,
		constants.CouponFilePath3,
	}
}

func loadEnvFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			_ = os.Setenv(key, value)
		}
	}
}
