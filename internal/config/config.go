package config

import "os"

type Config struct {
	DBPassword string
	DBUser     string
	DBName     string
}

func LoadConfig() (*Config, error) {
	return &Config{
		DBPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		DBUser:     getEnv("POSTGRES_USER", "postgres"),
		DBName:     getEnv("POSTGRES_DB", "postgres"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
