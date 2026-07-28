package platform

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort     string
	DatabaseURL string
}

func LoadConfig() (Config, error) {
	config := Config{
		AppPort:     os.Getenv("APP_PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	if config.AppPort == "" {
		config.AppPort = "8080"
	}

	if config.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return config, nil
}
