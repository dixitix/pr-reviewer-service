// Package config содержит логику загрузки конфигурации приложения
// из переменных окружения.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// HTTPConfig описывает настройки HTTP-сервера.
type HTTPConfig struct {
	Addr string
}

// DBConfig описывает настройки подключения к базе данных.
type DBConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// Config агрегирует все настройки приложения.
type Config struct {
	HTTP HTTPConfig
	DB   DBConfig
}

// Load загружает конфигурацию из переменных окружения и проверяет обязательные поля.
func Load() (Config, error) {
	cfg := Config{
		HTTP: HTTPConfig{
			Addr: getEnv("HTTP_ADDR", ":8080"),
		},
		DB: DBConfig{
			DSN:             os.Getenv("DATABASE_DSN"),
			MaxOpenConns:    mustParseInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:    mustParseInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: mustParseDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
	}

	if cfg.DB.DSN == "" {
		return Config{}, fmt.Errorf("DATABASE_DSN is required")
	}

	return cfg, nil
}

// getEnv возвращает значение переменной окружения или дефолт.
func getEnv(name, def string) string {
	if val, ok := os.LookupEnv(name); ok && val != "" {
		return val
	}

	return def
}

// mustParseInt парсит целое из переменной окружения или возвращает дефолт.
func mustParseInt(name string, def int) int {
	raw := os.Getenv(name)
	if raw == "" {
		return def
	}

	v, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}

	return v
}

// mustParseDuration парсит duration из переменной окружения или возвращает дефолт.
func mustParseDuration(name string, def time.Duration) time.Duration {
	raw := os.Getenv(name)
	if raw == "" {
		return def
	}

	d, err := time.ParseDuration(raw)
	if err != nil {
		return def
	}

	return d
}
