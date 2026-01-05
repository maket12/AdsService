package postgres

import (
	entity2 "ads/authservice/internal/domain/entity"
	"ads/authservice/pkg"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *pkg.Config) (*gorm.DB, error) {
	fmt.Printf("connecting to PostgreSQL... host=%s port=%d database=%s\n",
		cfg.DBHost, cfg.DBPort, cfg.DBName)

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName)

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = (DB.AutoMigrate(&entity2.User{}, &entity2.Session{}, &entity2.Profile{})); err != nil {
		return nil, err
	}

	fmt.Printf("✅ PostgreSQL connected successfully")
	return DB, nil
}
