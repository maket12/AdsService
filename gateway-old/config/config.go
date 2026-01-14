package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	AuthGrpcAddr  string `env:"AUTH_GRPC_ADDR" envDefault:"authservice:50051"`
	UserGrpcAddr  string `env:"USER_GRPC_ADDR" envDefault:"userservice:50052"`
	AdminGrpcAddr string `env:"ADMIN_GRPC_ADDR" envDefault:"adminservice:50053"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	// Валидация обязательных полей
	if cfg.AuthGrpcAddr == "" {
		return nil, fmt.Errorf("AUTH_GRPC_ADDR is required")
	}
	if cfg.UserGrpcAddr == "" {
		return nil, fmt.Errorf("USER_GRPC_ADDR is required")
	}
	if cfg.AdminGrpcAddr == "" {
		return nil, fmt.Errorf("ADMIN_GRPC_ADDR is required")
	}

	fmt.Printf("Config loaded successfully\n")
	fmt.Printf("   AuthGrpcAddr: %s\n", cfg.AuthGrpcAddr)
	fmt.Printf("   UserGrpcAddr: %s\n", cfg.UserGrpcAddr)
	fmt.Printf("   AdminGrpcAddr: %s\n", cfg.AdminGrpcAddr)

	return cfg, nil
}
