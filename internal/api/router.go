package api

import (
	"database/sql"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/michaelwp/student_attendance/internal/config"
	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/michaelwp/student_attendance/docs"
	"github.com/michaelwp/student_attendance/internal/api/handlers"
	"github.com/michaelwp/student_attendance/internal/repository"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func SetupRoutes(
	app *fiber.App,
	postgresClient *sql.DB,
	s3Client *s3.Client,
	s3Config *config.S3Config,
	redisClient *redis.Client,
) {

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize repositories and handlers
	repos := repository.NewRepositories(postgresClient)
	h := handlers.NewHandlers(&handlers.HandlerDependencies{
		Repositories: repos,
		S3Client:     s3Client,
		S3Config:     s3Config,
		RedisClient:  redisClient,
	})

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Health check route
	// @Summary Health check
	// @Description Check if the API is running
	// @Tags Health
	// @Produce json
	// @Success 200 {object} map[string]interface{} "API is healthy"
	// @Router /health [get]
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "OK",
			"message": "Student Attendance API is running",
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// Teacher routes
	teachers := api.Group("/teachers")
	teachers.Post("/", h.Teacher.Create)
	teachers.Get("/", h.Teacher.GetAll)
	teachers.Get("/:id", h.Teacher.GetByID)
	teachers.Get("/teacher-id/:teacherId", h.Teacher.GetByTeacherID)
	teachers.Put("/:id", h.Teacher.Update)
	teachers.Delete("/:id", h.Teacher.Delete)
	teachers.Put("/:id/photo", h.Teacher.UploadPhoto)
	teachers.Get("/:id/photo", h.Teacher.GetPhoto)
	teachers.Put("/teacher-id/:teacherId/reset-password", h.Student.ResetPassword)
	teachers.Put("/teacher-id/:teacherId/password", h.Student.UpdatePassword)

	// Class routes
	classes := api.Group("/classes")
	classes.Post("/", h.Class.Create)
	classes.Get("/", h.Class.GetAll)
	classes.Get("/:id", h.Class.GetByID)
	classes.Get("/teacher-id/:teacherId", h.Class.GetByTeacher)
	classes.Put("/:id", h.Class.Update)
	classes.Delete("/:id", h.Class.Delete)

	// Student routes
	students := api.Group("/students")
	students.Post("/", h.Student.Create)
	students.Get("/", h.Student.GetAll)
	students.Get("/:id", h.Student.GetByID)
	students.Get("/student-id/:studentId", h.Student.GetByStudentID)
	students.Get("/class-id/:classId", h.Student.GetByClass)
	students.Put("/:id", h.Student.Update)
	students.Delete("/:id", h.Student.Delete)
	students.Put("/:id/photo", h.Student.UploadPhoto)
	students.Get("/:id/photo", h.Student.GetPhoto)
	students.Put("/student-id/:studentId/reset-password", h.Student.ResetPassword)
	students.Put("/student-id/:studentId/password", h.Student.UpdatePassword)

	// Attendance routes
	attendances := api.Group("/attendances")
	attendances.Post("/", h.Attendance.Create)
	attendances.Get("/:id", h.Attendance.GetByID)
	attendances.Get("/student-id/:studentId", h.Attendance.GetByStudent)
	attendances.Get("/class-id/:classId", h.Attendance.GetByClass)
	attendances.Get("/date-range", h.Attendance.GetByDateRange)
	attendances.Put("/:id", h.Attendance.Update)
	attendances.Delete("/:id", h.Attendance.Delete)

	// Absent Request routes
	absentRequests := api.Group("/absent-requests")
	absentRequests.Post("/", h.AbsentRequest.Create)
	absentRequests.Get("/:id", h.AbsentRequest.GetByID)
	absentRequests.Get("/student-id/:studentId", h.AbsentRequest.GetByStudent)
	absentRequests.Get("/class-id/:classId", h.AbsentRequest.GetByClass)
	absentRequests.Get("/pending", h.AbsentRequest.GetPending)
	absentRequests.Patch("/:id/status", h.AbsentRequest.UpdateStatus)
	absentRequests.Delete("/:id", h.AbsentRequest.Delete)
}
