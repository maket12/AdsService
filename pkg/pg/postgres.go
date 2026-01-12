package pg

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresConfig struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string

	OpenConn     int
	IdleConn     int
	ConnLifeTime time.Duration
}

func NewPostgresConfig(
	host string, port int, user, password, name, ssl string,
	openConn, idleConn int, connLifeTime time.Duration) *PostgresConfig {
	return &PostgresConfig{
		DBHost:       host,
		DBPort:       port,
		DBUser:       user,
		DBPassword:   password,
		DBName:       name,
		SSLMode:      ssl,
		OpenConn:     openConn,
		IdleConn:     idleConn,
		ConnLifeTime: connLifeTime,
	}
}

func (pc *PostgresConfig) dsn() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pc.DBHost, pc.DBPort, pc.DBUser, pc.DBPassword, pc.DBName, pc.SSLMode,
	)
}

type PostgresClient struct {
	DB *sql.DB
}

func NewPostgresClient(config *PostgresConfig) (*PostgresClient, error) {
	var dsn = config.dsn()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	if config != nil {
		db.SetMaxOpenConns(config.OpenConn)
		db.SetMaxIdleConns(config.IdleConn)
		db.SetConnMaxLifetime(config.ConnLifeTime)
	} else {
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(time.Minute * 5)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresClient{DB: db}, nil
}

func (c *PostgresClient) Close() error {
	return c.DB.Close()
}
