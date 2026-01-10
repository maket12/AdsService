package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	// Database
	DBHost     string `env:"DB_HOST,required"`
	DBPort     int    `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	DBSSLMode  string `env:"DB_SSL_MODE" envDefault:"prefer"`

	DBOpenConn     int           `env:"DB_DB_OPEN_CONNECTIONS" envDefault:"25"`
	DBIdleConn     int           `env:"DB_DB_IDLE_CONNECTIONS" envDefault:"25"`
	DBConnLifeTime time.Duration `env:"DB_CONNECTION_LIFETIME" envDefault:"5m"`

	// JWT
	AccessSecret  string        `env:"ACCESS_SECRET,required"`
	RefreshSecret string        `env:"REFRESH_SECRET,required"`
	AccessTTL     time.Duration `env:"ACCESS_TTL" envDefault:"15m"`
	RefreshTTL    time.Duration `env:"REFRESH_TTL" envDefault:"720h"`

	// Password hasher
	PasswordCost int `env:"PASSWORD_COST" envDefault:"4"`

	// Service
	GRPCPort    int    `env:"GRPC_PORT" envDefault:"50051"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"INFO"`
	Environment string `env:"ENVIRONMENT" envDefault:"development"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	if cfg.DBHost == "" {
		return nil, fmt.Errorf("DB_HOST is required")
	}
	if cfg.DBUser == "" {
		return nil, fmt.Errorf("DBUser is required")
	}
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DBPassword is required")
	}
	if cfg.DBName == "" {
		return nil, fmt.Errorf("DBName is required")
	}

	if cfg.AccessSecret == "" {
		return nil, fmt.Errorf("AccessSecret is required")
	}
	if cfg.RefreshSecret == "" {
		return nil, fmt.Errorf("RefreshSecret is required")
	}

	fmt.Printf("Config loaded successfully\n")
	fmt.Printf("   Environment: %s\n", cfg.Environment)
	fmt.Printf("   Log Level: %s\n", cfg.LogLevel)
	fmt.Printf("   DB Host: %s\n", cfg.DBHost)
	fmt.Printf("   gRPC Port: %d\n", cfg.GRPCPort)

	return cfg, nil
}
