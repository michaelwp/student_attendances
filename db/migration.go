package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/michaelwp/student_attendance/internal/config"
	"log"
	"os"
)

func RunMigration(direction string) error {
	cfg := config.NewPostgresConfig()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	m, err := migrate.New("file://db/migrations", dsn)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %v", err)
	}
	defer m.Close()

	switch direction {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply up migrations: %v", err)
		}
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply down migrations: %v", err)
		}
	default:
		return fmt.Errorf("invalid migration direction: %s", direction)
	}

	return nil
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Migration direction (up/down) is required")
	}

	direction := os.Args[1]
	if err := RunMigration(direction); err != nil {
		log.Fatal(err)
	}

	log.Printf("Migration '%s' completed successfully", direction)
}
