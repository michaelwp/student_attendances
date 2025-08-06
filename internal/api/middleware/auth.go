package middleware

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelwp/student_attendance/pkg"
	"github.com/redis/go-redis/v9"
)

// JWTMiddleware validates JWT tokens and checks Redis cache
func JWTMiddleware(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.token_required",
				"error":         "Authorization token is required",
			})
		}

		// Check if the token starts with "Bearer"
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.invalid_token_format",
				"error":         "Token must be in format: Bearer <token>",
			})
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.token_required",
				"error":         "Authorization token is required",
			})
		}

		// Validate JWT token
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"translate_key": "error.jwt_secret_missing",
				"error":         "JWT secret not configured",
			})
		}

		jwtConfig := pkg.JWTConfig{
			SecretKey:     jwtSecret,
			TokenDuration: time.Hour,
		}

		claims, err := pkg.ValidateToken(token, jwtConfig)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.invalid_token",
				"error":         "Invalid or expired token",
			})
		}

		// Check if a token exists in Redis
		tokenKey := "token:" + claims.UserType + ":" + claims.UserID
		cachedToken, err := redisClient.Get(context.Background(), tokenKey).Result()
		if err != nil || cachedToken != token {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.token_not_found",
				"error":         "Token not found or expired",
			})
		}

		// Check if the token is expired
		if time.Now().After(claims.ExpiresAt.Time) {
			// Remove expired token from Redis
			redisClient.Del(context.Background(), tokenKey)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.token_expired",
				"error":         "Token has expired",
			})
		}

		// Store user information in context for use in handlers
		c.Locals("userID", claims.UserID)
		c.Locals("userType", claims.UserType)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// OptionalJWTMiddleware validates JWT tokens if present but doesn't require them
func OptionalJWTMiddleware(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// No token provided, continue without authentication
			return c.Next()
		}

		// Token provided, validate it
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.invalid_token_format",
				"error":         "Token must be in format: Bearer <token>",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return c.Next()
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			return c.Next()
		}

		jwtConfig := pkg.JWTConfig{
			SecretKey:     jwtSecret,
			TokenDuration: time.Hour,
		}

		claims, err := pkg.ValidateToken(token, jwtConfig)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.invalid_token",
				"error":         "Invalid or expired token",
			})
		}

		// Check if a token exists in Redis
		tokenKey := "token:" + claims.UserType + ":" + claims.UserID
		cachedToken, err := redisClient.Get(context.Background(), tokenKey).Result()
		if err != nil || cachedToken != token {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.token_not_found",
				"error":         "Token not found or expired",
			})
		}

		// Check if the token is expired
		if time.Now().After(claims.ExpiresAt.Time) {
			redisClient.Del(context.Background(), tokenKey)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.token_expired",
				"error":         "Token has expired",
			})
		}

		// Store user information in context
		c.Locals("userID", claims.UserID)
		c.Locals("userType", claims.UserType)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// RequireUserType middleware requires specific user types
func RequireUserType(allowedTypes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("userType")
		if userType == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"translate_key": "error.authentication_required",
				"error":         "Authentication required",
			})
		}

		userTypeStr := userType.(string)
		for _, allowedType := range allowedTypes {
			if userTypeStr == allowedType {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"translate_key": "error.insufficient_permissions",
			"error":         "Insufficient permissions for this operation",
		})
	}
}
