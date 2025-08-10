package api

import (
	"database/sql"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/michaelwp/student_attendance/internal/config"
	"github.com/michaelwp/student_attendance/internal/models"
	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/michaelwp/student_attendance/docs"
	"github.com/michaelwp/student_attendance/internal/api/handlers"
	"github.com/michaelwp/student_attendance/internal/api/middleware"
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
	teachers := api.Group("/teachers", middleware.JWTMiddleware(redisClient))
	teachers.Post("/", h.Teacher.Create)
	teachers.Get("/all", h.Teacher.GetAll)
	teachers.Get("/record-id/:id", h.Teacher.GetByID)
	teachers.Get("/teacher-id/:teacherId", h.Teacher.GetByTeacherID)
	teachers.Put("/record-id/:id", h.Teacher.Update)
	teachers.Delete("/record-id/:id", h.Teacher.Delete)
	teachers.Put("/record-id/:id/photo", h.Teacher.UploadPhoto)
	teachers.Get("/record-id/:id/photo", h.Teacher.GetPhoto)
	teachers.Put("/teacher-id/:teacherId/reset-password", h.Teacher.ResetPassword)
	teachers.Put("/teacher-id/:teacherId/password", h.Teacher.UpdatePassword)
	teachers.Get("/stats", h.Admin.GetStat)

	// Class routes
	classes := api.Group("/classes", middleware.JWTMiddleware(redisClient))
	classes.Post("/", h.Class.Create)
	classes.Get("/", h.Class.GetAll)
	classes.Get("/:id", h.Class.GetByID)
	classes.Get("/teacher-id/:teacherId", h.Class.GetByTeacher)
	classes.Put("/:id", h.Class.Update)
	classes.Delete("/:id", h.Class.Delete)

	// Student routes
	students := api.Group("/students", middleware.JWTMiddleware(redisClient))
	students.Post("/", h.Student.Create)
	students.Get("/all", h.Student.GetAll)
	students.Get("/record-id/:id", h.Student.GetByID)
	students.Get("/student-id/:studentId", h.Student.GetByStudentID)
	students.Get("/class-id/:classId", h.Student.GetByClass)
	students.Put("/record-id/:id", h.Student.Update)
	students.Delete("/record-id/:id", h.Student.Delete)
	students.Put("/record-id/:id/photo", h.Student.UploadPhoto)
	students.Get("/record-id/:id/photo", h.Student.GetPhoto)
	students.Put("/student-id/:studentId/reset-password", h.Student.ResetPassword)
	students.Put("/student-id/:studentId/password", h.Student.UpdatePassword)
	students.Get("/stats", h.Admin.GetStat)

	// Attendance routes
	attendances := api.Group("/attendances", middleware.JWTMiddleware(redisClient))
	attendances.Post("/", h.Attendance.Create)
	attendances.Get("/all", h.Attendance.GetAll)
	attendances.Get("/attendances-id/:id", h.Attendance.GetByID)
	attendances.Get("/student-id/:studentId", h.Attendance.GetByStudent)
	attendances.Get("/class-id/:classId", h.Attendance.GetByClass)
	attendances.Get("/date-range", h.Attendance.GetByDateRange)
	attendances.Put("/attendances-id/:id", h.Attendance.Update)
	attendances.Delete("/attendances-id/:id", h.Attendance.Delete)

	// Absent Request routes
	absentRequests := api.Group("/absent-requests",
		middleware.JWTMiddleware(redisClient),
		middleware.RequireUserType(
			models.UserTypeStudent.String(),
			models.UserTypeTeacher.String()),
	)
	absentRequests.Post("/", h.AbsentRequest.Create)
	absentRequests.Get("/absent-request-id/:id", h.AbsentRequest.GetByID)
	absentRequests.Get("/student-id/:studentId", h.AbsentRequest.GetByStudent)
	absentRequests.Get("/class-id/:classId", h.AbsentRequest.GetByClass)
	absentRequests.Get("/absent-request-id/pending", h.AbsentRequest.GetPending)
 absentRequests.Patch("/absent-request-id/:id/status", h.AbsentRequest.UpdateStatus)
	absentRequests.Put("/absent-request-id/:id", h.AbsentRequest.UpdateByCurrentStudent)
	absentRequests.Delete("/absent-request-id/:id", h.AbsentRequest.Delete)
	absentRequests.Get("/current-student", h.AbsentRequest.GetByCurrentStudent)

	// Admin routes
	admins := api.Group("/admins",
		middleware.JWTMiddleware(redisClient),
		middleware.RequireUserType(models.UserTypeAdmin.String()),
	)
	admins.Post("/", h.Admin.Create)
	admins.Put("/admin-id/:id", h.Admin.Update)
	admins.Delete("/admin-id/:id", h.Admin.Delete)
	admins.Get("/all", h.Admin.GetAll)
	admins.Get("/admin-id/:id", h.Admin.GetByID)
	admins.Get("/stats", h.Admin.GetStat)
	admins.Get("/email/:email", h.Admin.GetByEmail)
	admins.Put("/password", h.Admin.UpdatePassword)
	admins.Put("/admin-id/:id/status", h.Admin.SetActiveStatus)
	admins.Put("/admin-id/:id/reset-password", h.Admin.ResetPassword)

	// Authentication routes
	auth := api.Group("/auth")
	auth.Post("/login", h.Auth.Login)
	auth.Post("/logout", middleware.JWTMiddleware(redisClient), h.Auth.Logout)

	// Student dashboard routes (student authentication required)
	studentDashboard := api.Group("/student", middleware.JWTMiddleware(redisClient), middleware.RequireUserType(models.UserTypeStudent.String()))
	studentDashboard.Get("/profile", h.Student.GetProfile)
	studentDashboard.Put("/profile", h.Student.UpdateProfile)
 studentDashboard.Put("/password", h.Student.UpdateCurrentPassword)
	
	// Public student attendance marking (no authentication required)
	api.Post("/attendance/mark", h.Attendance.MarkAttendance)
}
