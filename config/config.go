package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTP     HTTPConfig
	Postgres PostgresConfig
	Storage  StorageConfig
	JWT      JWTConfig
	LogLevel string
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

type HTTPConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	MaxConns int32
}

type StorageConfig struct {
	BasePath      string
	MaxFileSizeMB int64
}

func New() (*Config, error) {
	ttl, err := strconv.Atoi(getEnv("JWT_TTL_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_TTL_HOURS: %w", err)
	}
	pgMaxConns, err := strconv.Atoi(getEnv("POSTGRES_MAX_CONNS", "10"))
	if err != nil {
		return nil, fmt.Errorf("invalid POSTGRES_MAX_CONNS: %w", err)
	}

	maxFileSizeMB, err := strconv.ParseInt(getEnv("STORAGE_MAX_FILE_SIZE_MB", "100"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid STORAGE_MAX_FILE_SIZE_MB: %w", err)
	}

	cfg := &Config{
		LogLevel: getEnv("LOG_LEVEL", "info"),
		HTTP: HTTPConfig{
			Host:         getEnv("HTTP_HOST", "134.17.130.253"),
			Port:         getEnv("HTTP_PORT", "8080"),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Postgres: PostgresConfig{
			Host:     requireEnv("POSTGRES_HOST"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     requireEnv("POSTGRES_USER"),
			Password: requireEnv("POSTGRES_PASSWORD"),
			DBName:   requireEnv("POSTGRES_DB"),
			MaxConns: int32(pgMaxConns),
		},
		Storage: StorageConfig{
			BasePath:      getEnv("STORAGE_BASE_PATH", "./data/files"),
			MaxFileSizeMB: maxFileSizeMB,
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "secret"),
			TTL:    time.Duration(ttl) * time.Hour,
		},
	}

	var errs []error
	if cfg.Postgres.Host == "" {
		errs = append(errs, errors.New("POSTGRES_HOST is required"))
	}
	if cfg.Postgres.User == "" {
		errs = append(errs, errors.New("POSTGRES_USER is required"))
	}
	if cfg.Postgres.Password == "" {
		errs = append(errs, errors.New("POSTGRES_PASSWORD is required"))
	}
	if cfg.Postgres.DBName == "" {
		errs = append(errs, errors.New("POSTGRES_DB is required"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func requireEnv(key string) string {
	return os.Getenv(key)
}
