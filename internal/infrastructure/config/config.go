package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Environment string `env:"ENV" envDefault:"development"`
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
	Cache       CacheConfig
	Cors        CorsConfig
}

type ServerConfig struct {
	Port         int           `env:"SERVER_PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" envDefault:"60s"`
}

type DatabaseConfig struct {
	Path string `env:"DB_PATH" envDefault:"./data.db"`
}

type JWTConfig struct {
	Secret           string        `env:"JWT_SECRET" envDefault:"your-secret-key"`
	ExpirationHours  time.Duration `env:"JWT_EXPIRATION_HOURS" envDefault:"24h"`
	RefreshDuration  time.Duration `env:"JWT_REFRESH_DURATION" envDefault:"168h"`
	SigningAlgorithm string        `env:"JWT_SIGNING_ALGORITHM" envDefault:"HS256"`
}

type CacheConfig struct {
	DefaultExpiration time.Duration `env:"CACHE_DEFAULT_EXPIRATION" envDefault:"5m"`
	CleanupInterval   time.Duration `env:"CACHE_CLEANUP_INTERVAL" envDefault:"10m"`
}

type CorsConfig struct {
	AllowedOrigins   []string `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`
	AllowedMethods   []string `env:"CORS_ALLOWED_METHODS" envDefault:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowedHeaders   []string `env:"CORS_ALLOWED_HEADERS" envDefault:"Accept,Authorization,Content-Type,X-CSRF-Token"`
	ExposedHeaders   []string `env:"CORS_EXPOSED_HEADERS" envDefault:"Link"`
	AllowCredentials bool     `env:"CORS_ALLOW_CREDENTIALS" envDefault:"true"`
	MaxAge           int      `env:"CORS_MAX_AGE" envDefault:"300"`
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found")
	}

	config := &Config{
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Database: DatabaseConfig{
			Path: "./data.db",
		},
		JWT: JWTConfig{
			Secret:           "your-secret-key",
			ExpirationHours:  24 * time.Hour,
			RefreshDuration:  168 * time.Hour,
			SigningAlgorithm: "HS256",
		},
		Cache: CacheConfig{
			DefaultExpiration: 5 * time.Minute,
			CleanupInterval:   10 * time.Minute,
		},
		Cors: CorsConfig{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		},
	}

	return config, nil
}