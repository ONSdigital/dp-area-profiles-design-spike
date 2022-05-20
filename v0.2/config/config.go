package config

import (
	"fmt"
	"os"
)

type Config struct {
	Username string
	Password string
	Database string
}

// Get return the app config.
func Get() (*Config, error) {
	dbName := os.Getenv("AP_DATABASE_NAME")
	if dbName == "" {
		return nil, fmt.Errorf("expected env var %q but not found", "AP_DATABASE_NAME")
	}

	dbUsername := os.Getenv("AP_POSTGRES_USER")
	if dbUsername == "" {
		return nil, fmt.Errorf("expected env var %q but not found", "AP_DATABASE_USERNAME")
	}

	dbPassword := os.Getenv("AP_POSTGRES_PASSWORD")
	if dbPassword == "" {
		return nil, fmt.Errorf("expected env var %q but not found", "AP_DATABASE_PASSWORD")
	}

	return &Config{
		Username: dbUsername,
		Password: dbPassword,
		Database: dbName,
	}, nil
}
