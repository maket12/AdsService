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

	// JWT
	AccessSecret  string        `env:"ACCESS_SECRET,required"`
	RefreshSecret string        `env:"REFRESH_SECRET,required"`
	AccessTTL     time.Duration `env:"ACCESS_TTL" envDefault:"15m"`
	RefreshTTL    time.Duration `env:"REFRESH_TTL" envDefault:"720h"`

	// Password hasher
	PasswordCost int `env:"PASSWORD_COST" envDefault:"4"`

	// RabbitMQ
	RabbitHost     string `env:"RABBIT_HOST,required"`
	RabbitPort     int    `env:"RABBIT_PORT" envDefault:"5672"`
	RabbitUser     string `env:"RABBIT_USER,required"`
	RabbitPassword string `env:"RABBIT_PASSWORD,required"`

	RabbitWaitTime time.Duration `env:"RABBIT_WAIT_TIME" envDefault:"30s"`
	RabbitAttempts int           `env:"RABBIT_ATTEMPTS" envDefault:"5"`

	ExchangeName string `env:"EXCHANGE_NAME" envDefault:"account_topic"`
	RoutingKey   string `env:"ROUTING_KEY,required"`

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
		return nil, fmt.Errorf("PgUser is required")
	}
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("PgPassword is required")
	}
	if cfg.DBName == "" {
		return nil, fmt.Errorf("PgDBName is required")
	}

	if cfg.AccessSecret == "" {
		return nil, fmt.Errorf("AccessSecret is required")
	}
	if cfg.RefreshSecret == "" {
		return nil, fmt.Errorf("RefreshSecret is required")
	}

	if cfg.RabbitHost == "" {
		return nil, fmt.Errorf("RabbitUser is required")
	}
	if cfg.RabbitUser == "" {
		return nil, fmt.Errorf("RabbitUser is required")
	}
	if cfg.RabbitPassword == "" {
		return nil, fmt.Errorf("RabbitPassword is required")
	}

	if cfg.RoutingKey == "" {
		return nil, fmt.Errorf("RoutingKey is required")
	}

	fmt.Printf("Config loaded successfully\n")
	fmt.Printf("   Environment: %s\n", cfg.Environment)
	fmt.Printf("   Log Level: %s\n", cfg.LogLevel)
	fmt.Printf("   DB Host: %s\n", cfg.DBHost)
	fmt.Printf("   RabbitMQ Host: %s\n", cfg.RabbitHost)
	fmt.Printf("   gRPC Port: %d\n", cfg.GRPCPort)

	return cfg, nil
}
