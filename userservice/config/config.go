package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	// Database
	DBHost     string `env:"DB_HOST,required"`
	DBPort     int    `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`

	// MongoDB
	MongoURI    string `env:"MONGODB_URI,required"`
	MongoDB     string `env:"MONGODB_DB_NAME" envDefault:"ads_service"`
	MongoBucket string `env:"MONGODB_BUCKET_NAME" envDefault:"photos"`

	// JWT
	JWTAccessSecret  string `env:"JWT_ACCESS_SECRET,required"`
	JWTRefreshSecret string `env:"JWT_REFRESH_SECRET,required"`

	// Service
	GRPCPort    int    `env:"GRPC_PORT" envDefault:"50052"`
	Environment string `env:"ENVIRONMENT" envDefault:"development"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	// Валидация обязательных полей
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

	if cfg.MongoURI == "" {
		return nil, fmt.Errorf("MongoURI is required")
	}

	if cfg.JWTAccessSecret == "" {
		return nil, fmt.Errorf("JWT_Access_SECRET is required")
	}
	if cfg.JWTRefreshSecret == "" {
		return nil, fmt.Errorf("JWTRefreshSecret is required")
	}

	fmt.Printf("Config loaded successfully\n")
	fmt.Printf("   Environment: %s\n", cfg.Environment)
	fmt.Printf("   DB Host: %s\n", cfg.DBHost)
	fmt.Printf("   gRPC Port: %d\n", cfg.GRPCPort)

	return cfg, nil
}
