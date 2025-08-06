// Student Attendance API
//
// # A comprehensive API for managing student attendance, classes, teachers, and absence requests
//
// Terms Of Service: N/A
// Schemes: http, https
// Host: localhost:8080
// BasePath: /api/v1
// Version: 1.0.0
// License: MIT https://opensource.org/licenses/MIT
// Contact: API Support <support@studentattendance.com>
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/michaelwp/student_attendance/docs"
	"github.com/michaelwp/student_attendance/internal/api"
	"github.com/michaelwp/student_attendance/internal/config"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	if os.Getenv("ENVIRONMENT") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}

func gracefulShutdown(app *fiber.App, postgresClient *sql.DB, postgresConfig *config.PostgresConfig, redisClient *redis.Client) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if err := postgresConfig.CloseDB(postgresClient); err != nil {
		log.Printf("Error closing PostgreSQL connection: %v\n", err)
	}

	if err := config.CloseRedisConnection(redisClient); err != nil {
		log.Printf("Error closing Redis connection: %v\n", err)
	}

	if err := app.Shutdown(); err != nil {
		log.Printf("Error shutting down Fiber server: %v\n", err)
	}

	log.Println("Server gracefully stopped")
}

func main() {
	// connect to PostgreSQL database
	postgresConfig := config.NewPostgresConfig()
	postgresClient, err := postgresConfig.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to postgres database: %v", err)
	}

	// connect to AWS S3
	s3Config := config.NewS3Config()
	s3Client, err := s3Config.NewS3Client()
	if err != nil {
		log.Fatalf("Error connecting to AWS S3: %v", err)
	}

	// connect to Redis
	redisClient, err := config.NewRedisClient()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		EnablePrintRoutes:     true,
		AppName:               "Student Attendance API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logLevel := os.Getenv("LOG_LEVEL")
			if logLevel == "debug" {
				log.Printf("[ERROR] %v\n", err)
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"translate.key": "internal_server_error",
				"error":         "Internal Server Error",
			})
		},
	})

	// Setup routes
	api.SetupRoutes(app, postgresClient, s3Client, s3Config, redisClient)

	port := os.Getenv("PORT")

	go func() {
		// Wait for a shutdown signal
		gracefulShutdown(app, postgresClient, postgresConfig, redisClient)
	}()

	log.Printf("Student Attendance API listening on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
