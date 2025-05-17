package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
)

type Config struct {
	Database databaseConfig
	Server   serverConfig
}
type serverConfig struct {
	Port string
	Host string
}
type databaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func New() Config {
	return Config{
		Server: serverConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Database: databaseConfig{
			Host:     getEnv("DATABASE_HOST", "db"),
			Port:     getEnv("DATABASE_PORT", "5432"),
			User:     getEnv("DATABASE_USER", "postgres"),
			Password: getEnv("DATABASE_PASSWORD", "postgres"),
			DBName:   getEnv("DATABASE_DB", "htmxtst"),
			SSLMode:  getEnv("DATABASE_SSL_MODE", "disable"),
		},
	}
}

func (cfg *databaseConfig) DatabaseURL() string {
	encodedPassword := url.QueryEscape(cfg.Password)

	hostPort := net.JoinHostPort(cfg.Host, cfg.Port)

	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.User,
		encodedPassword,
		hostPort,
		cfg.DBName,
		cfg.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}

	return defaultValue
}
