package config

import "os"

type Config struct {
	EngineURL string
	Port      string
	GinMode   string
}

func Load() *Config {
	return &Config{
		EngineURL: getEnv("ENGINE_URL", "http://localhost:8000"),
		Port:      getEnv("PORT", "8080"),
		GinMode:   getEnv("GIN_MODE", "debug"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}