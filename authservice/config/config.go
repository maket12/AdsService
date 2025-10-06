package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"log"
)

type Config struct {
	// Database
	DBHost     string `env:"DB_HOST,required"`
	DBPort     int    `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`

	// JWT
	JWTAccessSecret  string `env:"JWT_ACCESS_SECRET,required"`
	JWTRefreshSecret string `env:"JWT_REFRESH_SECRET,required"`

	// Service
	GRPCPort    int    `env:"GRPC_PORT" envDefault:"50051"`
	Environment string `env:"ENVIRONMENT" envDefault:"development"`
}

func Load() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Валидация обязательных полей
	if cfg.DBHost == "" {
		log.Fatal("DB_HOST is required")
	}
	if cfg.DBUser == "" {
		log.Fatal("DBUser is required")
	}
	if cfg.DBPassword == "" {
		log.Fatal("DBPassword is required")
	}
	if cfg.DBName == "" {
		log.Fatal("DBName is required")
	}

	if cfg.JWTAccessSecret == "" {
		log.Fatal("JWT_Access_SECRET is required")
	}
	if cfg.JWTRefreshSecret == "" {
		log.Fatal("JWTRefreshSecret is required")
	}

	fmt.Printf("Config loaded successfully\n")
	fmt.Printf("   Environment: %s\n", cfg.Environment)
	fmt.Printf("   DB Host: %s\n", cfg.DBHost)
	fmt.Printf("   gRPC Port: %d\n", cfg.GRPCPort)

	return cfg
}
