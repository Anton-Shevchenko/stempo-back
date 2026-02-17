package config

import (
	"os"
	"strconv"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
	JWT    JWTConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type JWTConfig struct {
	Secret string
}

func Load() *Config {
	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "stempo"),
			Password: getEnv("DB_PASSWORD", "stempo_password"),
			Name:     getEnv("DB_NAME", "stempo_db"),
		},
		Server: ServerConfig{
			Port:    getEnv("PORT", "3000"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "default-secret-key"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
