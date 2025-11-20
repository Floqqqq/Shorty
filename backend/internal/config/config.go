package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
    Port            string
    BaseURL         string
    PgHost          string
    PgPort          string
    PgUser          string
    PgPassword      string
    PgDB            string
    PgSSLMode       string
    ShortCodeLength int
}

func LoadConfigFromEnv() (*Config, error) {
    length := 6
    if s := os.Getenv("SHORT_CODE_LENGTH"); s != "" {
        if v, err := strconv.Atoi(s); err == nil {
            length = v
        }
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    baseURL := os.Getenv("BASE_URL")
    if baseURL == "" {
        return nil, errors.New("BASE_URL required")
    }
    return &Config{
        Port:           port,          // <- добавь сюда
        BaseURL:        baseURL,       // <- и базовый URL тоже
        PgHost:    os.Getenv("PGHOST"),
        PgPort:    os.Getenv("PGPORT"),
        PgUser:    os.Getenv("PGUSER"),
        PgPassword: os.Getenv("PGPASSWORD"),
        PgDB:      os.Getenv("PGDB"),
        PgSSLMode: os.Getenv("PGSSLMODE"),
        ShortCodeLength: length,
    }, nil
}

func envOr(k, def string) string {
    if v := os.Getenv(k); v != "" {
        return v
    }
    return def
}