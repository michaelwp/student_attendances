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
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/michaelwp/student_attendance/docs"
	"github.com/michaelwp/student_attendance/internal/api"
	"github.com/michaelwp/student_attendance/internal/config"
	"log"
	"os"
)

func init() {
	if os.Getenv("ENVIRONMENT") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}

func main() {
	// connect to PostgreSQL database
	postgres := config.NewPostgresConfig()
	postgresClient, err := postgres.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to postgres database: %v", err)
	}

	// connect to AWS S3
	s3Config := config.NewS3Config()
	s3Client, err := s3Config.NewS3Client()
	if err != nil {
		log.Fatalf("Error connecting to AWS S3: %v", err)
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
	api.SetupRoutes(app, postgresClient, s3Client, s3Config)

	port := os.Getenv("PORT")

	log.Printf("Student Attendance API listening on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
