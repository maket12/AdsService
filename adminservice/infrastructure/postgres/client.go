package postgres

import (
	"AdsService/adminservice/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"

	"AdsService/adminservice/domain/entity"
)

var DB *gorm.DB

func InitDB(cfg *config.Config, logger *slog.Logger) error {
	logger.Info("connecting to PostgreSQL...",
		slog.String("host", cfg.DBHost),
		slog.Int("port", cfg.DBPort),
		slog.String("database", cfg.DBName),
	)

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := (DB.AutoMigrate(&entity.User{}, &entity.Profile{})); err != nil {
		return err
	}

	logger.Info("Successfully connected to database!")

	return nil
}
