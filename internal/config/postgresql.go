package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}
}

func (c *PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func (c *PostgresConfig) ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", c.DSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err = db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	if err = c.SetupConnectionPool(db); err != nil {
		return nil, fmt.Errorf("error setting up connection pool: %v", err)
	}

	return db, nil
}

func (c *PostgresConfig) SetupConnectionPool(db *sql.DB) error {
	maxConn := os.Getenv("MAX_CONNECTIONS")
	maxIdleConn := os.Getenv("MAX_IDLE_CONNECTIONS")
	maxLifetimeConn := os.Getenv("MAX_LIFETIME_CONNECTIONS")

	if maxConn != "" {
		if maxCn, err := strconv.Atoi(maxConn); err == nil {
			db.SetMaxOpenConns(maxCn)
		}
	}

	if maxIdleConn != "" {
		if maxIdle, err := strconv.Atoi(maxIdleConn); err == nil {
			db.SetMaxIdleConns(maxIdle)
		}
	}

	if maxLifetimeConn != "" {
		if maxLife, err := strconv.Atoi(maxLifetimeConn); err == nil {
			db.SetConnMaxLifetime(time.Duration(maxLife) * time.Second)
		}
	}

	return nil
}

func (c *PostgresConfig) CloseDB(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("error closing database connection: %v", err)
	}
	return nil
}
