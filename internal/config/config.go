package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Environment string
	Port        int
	LogLevel    string
	DBHost      string
	DBPort      int
	DBName      string
	DBUser      string
	DBPassword  string
}

func Load() (*Config, error) {
	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "local"),
		Port:        getEnvAsInt("PORT", 8080),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnvAsInt("DB_PORT", 5432),
		DBName:      getEnv("DB_NAME", "school"),
		DBUser:      getEnv("DB_USER", ""),
		DBPassword:  getEnv("DB_PASSWORD", ""),
	}

	if cfg.DBHost == "" {
		return nil, fmt.Errorf("DB_HOST is required")
	}

	return cfg, nil
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
