package pkg

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plaintext password using bcrypt with specified cost.
// cost must be between bcrypt.MinCost (4) and bcrypt.MaxCost (31).
// Returns the hashed password as a string and any error encountered.
func HashPassword(password string, cost int) (string, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// ComparePasswords compares a bcrypt hashed password with a plaintext password.
// Returns nil on success, or an error on failure.
func ComparePasswords(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
