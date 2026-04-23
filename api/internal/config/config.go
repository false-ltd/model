package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Sync     SyncConfig
	Auth     AuthConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	dsn string
}

type SyncConfig struct {
	CooldownMinutes int
	ModelsDevURL    string
	CronMinutes     int
}

type AuthConfig struct {
	APIKeys []string
}

type CORSConfig struct {
	AllowedOrigins string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("MODEL_SERVER_PORT", "8080"),
			Mode: getEnv("MODEL_GIN_MODE", "release"),
		},
		Database: DatabaseConfig{
			dsn: getEnv("MODEL_DATABASE_DSN", "root:123456@tcp(127.0.0.1:3306)/models?charset=utf8mb4&parseTime=True&loc=Local"),
		},
		Sync: SyncConfig{
			CooldownMinutes: getEnvInt("MODEL_SYNC_COOLDOWN_MINUTES", 10),
			ModelsDevURL:    getEnv("MODEL_MODELS_DEV_URL", "https://models.dev/api.json"),
			CronMinutes:     getEnvInt("MODEL_SYNC_CRON_MINUTES", 60),
		},
		Auth: AuthConfig{
			APIKeys: getEnvSlice("MODEL_API_KEYS", []string{}),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("MODEL_CORS_ALLOWED_ORIGINS", "*"),
		},
	}
}

func (d *DatabaseConfig) DSN() string {
	return d.dsn
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvSlice(key string, fallback []string) []string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	var result []string
	start := 0
	for i := 0; i <= len(v); i++ {
		if i == len(v) || v[i] == ',' {
			part := v[start:i]
			if part != "" {
				result = append(result, part)
			}
			start = i + 1
		}
	}
	if len(result) == 0 {
		return fallback
	}
	return result
}
