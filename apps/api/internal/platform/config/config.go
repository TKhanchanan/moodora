package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultTimezone = "Asia/Bangkok"

type Config struct {
	AppName     string
	AppEnv      string
	AppPort     int
	AppTimezone *time.Location

	DatabaseURL string
	RedisURL    string
	Storage     StorageConfig
	DevUserID   string
	CORS        CORSConfig
}

type StorageConfig struct {
	Endpoint      string
	Region        string
	Bucket        string
	AccessKey     string
	SecretKey     string
	PublicBaseURL string
	UseSSL        bool
}

type CORSConfig struct {
	AllowedOrigins []string
}

func Load() (Config, error) {
	timezoneName := getEnv("APP_TIMEZONE", defaultTimezone)
	timezone, err := time.LoadLocation(timezoneName)
	if err != nil {
		return Config{}, fmt.Errorf("load timezone %q: %w", timezoneName, err)
	}

	port, err := strconv.Atoi(getEnv("APP_PORT", "8080"))
	if err != nil {
		return Config{}, fmt.Errorf("parse APP_PORT: %w", err)
	}

	cfg := Config{
		AppName:     getEnv("APP_NAME", "Moodora"),
		AppEnv:      getEnv("APP_ENV", "local"),
		AppPort:     port,
		AppTimezone: timezone,
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),
		DevUserID:   os.Getenv("DEV_USER_ID"),
		CORS: CORSConfig{
			AllowedOrigins: splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://127.0.0.1:3000")),
		},
		Storage: StorageConfig{
			Endpoint:      getEnv("S3_ENDPOINT", "http://localhost:9000"),
			Region:        getEnv("S3_REGION", "auto"),
			Bucket:        os.Getenv("S3_BUCKET"),
			AccessKey:     os.Getenv("S3_ACCESS_KEY"),
			SecretKey:     os.Getenv("S3_SECRET_KEY"),
			PublicBaseURL: os.Getenv("S3_PUBLIC_BASE_URL"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if c.AppPort <= 0 || c.AppPort > 65535 {
		return fmt.Errorf("APP_PORT must be between 1 and 65535")
	}
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}
	if c.RedisURL == "" {
		return errors.New("REDIS_URL is required")
	}
	if c.Storage.Bucket == "" {
		return errors.New("S3_BUCKET is required")
	}
	if c.Storage.AccessKey == "" {
		return errors.New("S3_ACCESS_KEY is required")
	}
	if c.Storage.SecretKey == "" {
		return errors.New("S3_SECRET_KEY is required")
	}
	return nil
}

func (c Config) HTTPAddr() string {
	return fmt.Sprintf(":%d", c.AppPort)
}

func (c Config) TimezoneName() string {
	if c.AppTimezone == nil {
		return defaultTimezone
	}
	return c.AppTimezone.String()
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			values = append(values, part)
		}
	}
	return values
}
