package platform

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	AppEnv          string
	AppPort         string
	DatabaseURL     string
	ShutdownTimeout time.Duration
}

func LoadConfig() (Config, error) {
	config := Config{
		AppEnv:          os.Getenv("APP_ENV"),
		AppPort:         os.Getenv("APP_PORT"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		ShutdownTimeout: 10 * time.Second,
	}

	if config.AppEnv == "" {
		config.AppEnv = "local"
	}

	if config.AppPort == "" {
		config.AppPort = "8080"
	}

	if config.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	if value := os.Getenv("SHUTDOWN_TIMEOUT"); value != "" {
		shutdownTimeout, err := time.ParseDuration(value)
		if err != nil {
			return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT is invalid: %w", err)
		}

		config.ShutdownTimeout = shutdownTimeout
	}

	return config, nil
}
