package postgres

import (
	"ads/userservice/domain/entity"
	"ads/userservice/pkg"
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

	if err = (DB.AutoMigrate(&entity.Profile{})); err != nil {
		return nil, err
	}

	fmt.Print("✅ PostgreSQL connected successfully")
	return DB, nil
}
