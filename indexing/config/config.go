package config

import (
	"errors"
	"os"
)

type Config struct {
	ServerAddr string
	DBDSN      string
}

func Load() (Config, error) {
	cfg := Config{
		ServerAddr: getEnv("SERVER_ADDR", ":8082"),
		DBDSN:      getEnv("DB_DSN", "host=localhost user=cfguser password=cfgpass dbname=cfgdb port=5432 sslmode=disable"),
	}

	if cfg.DBDSN == "" {
		return Config{}, errors.New("DB_DSN is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
