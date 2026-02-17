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

	DBOpenConn     int           `env:"DB_OPEN_CONNECTIONS" envDefault:"25"`
	DBIdleConn     int           `env:"DB_IDLE_CONNECTIONS" envDefault:"25"`
	DBConnLifeTime time.Duration `env:"DB_CONNECTION_LIFETIME" envDefault:"5m"`

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

	fmt.Printf("Config loaded successfully\n")
	fmt.Printf("   Environment: %s\n", cfg.Environment)
	fmt.Printf("   Log Level: %s\n", cfg.LogLevel)
	fmt.Printf("   DB Host: %s\n", cfg.DBHost)
	fmt.Printf("   gRPC Port: %d\n", cfg.GRPCPort)

	return cfg, nil
}
