package config

import (
	"os"
	"strconv"
)

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	JWTSecret    string
	Environment  string
	CORSOrigins  []string
	MaxPageSize  int
	DefaultLimit int
}

func LoadConfig() *Config {
	cfg := &Config{
		Environment: getEnv("ENV", "development"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		MaxPageSize: 100,
		DefaultLimit: 50,
	}

	// Database config
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.Name = getEnv("DB_NAME", "finanzas_app")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	// CORS
	corsOrigins := getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173")
	if corsOrigins != "" {
		cfg.CORSOrigins = []string{corsOrigins}
	}

	return cfg
}

func (c *Config) GetDSN() string {
	return "postgres://" + c.Database.User + ":" + c.Database.Password +
		"@" + c.Database.Host + ":" + c.Database.Port +
		"/" + c.Database.Name + "?sslmode=" + c.Database.SSLMode
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
