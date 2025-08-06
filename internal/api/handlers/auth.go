package handlers

import (
	"context"
	"github.com/michaelwp/student_attendance/internal/models"
	"log"
	"os"
	"strconv"
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

	var userID uint
	var storedPassword string
	var err error

	// Authenticate based on a user type
	switch loginReq.UserType {
	case models.UserTypeAdmin.String():
		userID, storedPassword, err = h.authenticateAdmin(c.Context(), loginReq.UserID)
	case models.UserTypeTeacher.String(), "teacher":
		userID, storedPassword, err = h.authenticateTeacher(c.Context(), loginReq.UserID)
	case models.UserTypeStudent.String():
		userID, storedPassword, err = h.authenticateStudent(c.Context(), loginReq.UserID)
	}

	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"translate_key": "error." + fiberErr.Message,
				"error":         fiberErr.Message,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.authentication_failed",
			"error":         "Authentication failed",
		})
	}

	// Verify password
	if err := pkg.ComparePasswords(storedPassword, loginReq.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"translate_key": "error.invalid_credentials",
			"error":         "Invalid credentials",
		})
	}

	// Generate and cache JWT token
	strUserID := strconv.FormatUint(uint64(userID), 10) // Safe conversion from uint to string
	token, err := h.generateAndCacheToken(c.Context(), strUserID, loginReq.UserType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"translate_key": "error.token_generation_failed",
			"error":         "Failed to generate or cache token",
		})
	}

	// Update last login for admin
	if loginReq.UserType == "admin" {
		h.updateAdminLastLogin(c.Context(), loginReq.UserID)
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

// authenticateAdmin validates admin credentials and returns userID and password if successful
func (h *authHandler) authenticateAdmin(ctx context.Context, userID string) (uint, string, error) {
	userExists, err := h.adminRepo.IsAdminExist(ctx, userID)
	if err != nil {
		return 0, "", err
	}
	if !userExists {
		return 0, "", fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Check if admin is active
	admin, err := h.adminRepo.GetByEmail(ctx, userID)
	if err != nil {
		return 0, "", err
	}

	if !admin.IsActive {
		return 0, "", fiber.NewError(fiber.StatusUnauthorized, "Account is deactivated")
	}

	storedPassword, err := h.adminRepo.GetPasswordByEmail(ctx, userID)
	if err != nil {
		return 0, "", err
	}

	return admin.ID, storedPassword, nil
}

// authenticateTeacher validates teacher credentials and returns userID and password if successful
func (h *authHandler) authenticateTeacher(ctx context.Context, userID string) (uint, string, error) {
	userExists, err := h.teacherRepo.IsTeacherExist(ctx, userID)
	if err != nil {
		return 0, "", err
	}
	if !userExists {
		return 0, "", fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Check if the teacher is active
	teacher, err := h.teacherRepo.GetByTeacherID(ctx, userID)
	if err != nil {
		return 0, "", err
	}

	if !teacher.IsActive {
		return 0, "", fiber.NewError(fiber.StatusUnauthorized, "Account is deactivated")
	}

	storedPassword, err := h.teacherRepo.GetPasswordByTeacherID(ctx, userID)
	if err != nil {
		return 0, "", err
	}

	return teacher.ID, storedPassword, nil
}

// authenticateStudent validates student credentials and returns userID and password if successful
func (h *authHandler) authenticateStudent(ctx context.Context, userID string) (uint, string, error) {
	userExists, err := h.studentRepo.IsStudentExist(ctx, userID)
	if err != nil {
		return 0, "", err
	}
	if !userExists {
		return 0, "", fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Check if the student is active
	student, err := h.studentRepo.GetByStudentID(ctx, userID)
	if err != nil {
		return 0, "", err
	}

	if !student.IsActive {
		return 0, "", fiber.NewError(fiber.StatusUnauthorized, "Account is deactivated")
	}

	storedPassword, err := h.studentRepo.GetPasswordByStudentID(ctx, userID)
	if err != nil {
		return 0, "", err
	}

	return student.ID, storedPassword, nil
}

// generateAndCacheToken generates a JWT token and caches it in Redis
func (h *authHandler) generateAndCacheToken(ctx context.Context, userID, userType string) (string, error) {
	// Get JWT secret from the environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fiber.NewError(fiber.StatusInternalServerError, "JWT secret not configured")
	}

	// Configure JWT
	jwtConfig := pkg.JWTConfig{
		SecretKey:     jwtSecret,
		TokenDuration: time.Hour, // 1-hour expiration
	}

	// Generate token
	token, err := pkg.GenerateToken(userID, userType, jwtConfig)
	if err != nil {
		return "", err
	}

	// Cache token in Redis with 1-hour expiration
	tokenKey := "token:" + userType + ":" + userID
	err = h.redisClient.Set(ctx, tokenKey, token, time.Hour).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

// updateAdminLastLogin updates the last login time for an admin
func (h *authHandler) updateAdminLastLogin(ctx context.Context, email string) {
	admin, err := h.adminRepo.GetByEmail(ctx, email)
	if err != nil || admin == nil {
		log.Println("error getting admin for last login update:", err)
		return
	}

	err = h.adminRepo.UpdateLastLogin(ctx, admin.ID, time.Now())
	if err != nil {
		log.Println("error updating last login for admin:", err)
	}
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
