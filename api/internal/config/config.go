package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Sync     SyncConfig
	Auth     AuthConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port           string
	Mode           string
	SiteURL        string
	TrustedProxies []string
}

type DatabaseConfig struct {
	dsn string
}

type SyncConfig struct {
	CooldownMinutes int
	ModelsDevURL    string
	CronMinutes     int
	// PublicTrigger allows unauthenticated POST /sync (cooldown, rate limit
	// and the sync mutex still apply) so the frontend refresh button works.
	PublicTrigger bool
}

type AuthConfig struct {
	APIKeys []string
}

type CORSConfig struct {
	AllowedOrigins string
}

func Load() (*Config, error) {
	dsn := os.Getenv("MODEL_DATABASE_DSN")
	if dsn == "" {
		return nil, errors.New("MODEL_DATABASE_DSN is required")
	}

	return &Config{
		Server: ServerConfig{
			Port:           getEnv("MODEL_SERVER_PORT", "8080"),
			Mode:           getEnv("MODEL_GIN_MODE", "release"),
			SiteURL:        getEnv("MODEL_SITE_URL", "https://model.false.ltd"),
			TrustedProxies: getEnvSlice("MODEL_TRUSTED_PROXIES", nil),
		},
		Database: DatabaseConfig{
			dsn: dsn,
		},
		Sync: SyncConfig{
			CooldownMinutes: getEnvInt("MODEL_SYNC_COOLDOWN_MINUTES", 10),
			ModelsDevURL:    getEnv("MODEL_MODELS_DEV_URL", "https://models.dev/api.json"),
			CronMinutes:     getEnvInt("MODEL_SYNC_CRON_MINUTES", 60),
			PublicTrigger:   getEnvBool("MODEL_SYNC_PUBLIC_TRIGGER", false),
		},
		Auth: AuthConfig{
			APIKeys: getEnvSlice("MODEL_API_KEYS", nil),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("MODEL_CORS_ALLOWED_ORIGINS", "*"),
		},
	}, nil
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

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		switch strings.ToLower(v) {
		case "1", "true", "yes":
			return true
		case "0", "false", "no":
			return false
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
	for _, part := range strings.Split(v, ",") {
		if part = strings.TrimSpace(part); part != "" {
			result = append(result, part)
		}
	}
	if len(result) == 0 {
		return fallback
	}
	return result
}
