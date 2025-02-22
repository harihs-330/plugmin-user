package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// LoadConfig loads the configuration for the application based on the given appName.
func Load(appName string) (*EnvConfig, error) {
	var cfg EnvConfig

	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	err := envconfig.Process(appName, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
