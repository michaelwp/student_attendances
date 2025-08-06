package pkg

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars    = "0123456789"
	symbolChars    = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

// GeneratePassword generates a random password with the specified length
// containing at least one symbol, number, uppercase, and lowercase character
func GeneratePassword(length int) (string, error) {
	minLength := 6 // Minimum length to accommodate all required character types

	if length < minLength {
		length = minLength
	}

	// Initialize with one character of each required type
	var password strings.Builder
	required := []string{
		symbolChars,
		numberChars,
		uppercaseChars,
		lowercaseChars,
	}

	// Add one character from each required set
	for _, charSet := range required {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "", err
		}
		password.WriteByte(charSet[n.Int64()])
	}

	// Fill the rest with random characters from all sets
	allChars := symbolChars + numberChars + uppercaseChars + lowercaseChars
	for i := minLength; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", err
		}
		password.WriteByte(allChars[n.Int64()])
	}

	// Convert to string and shuffle
	result := []rune(password.String())
	for i := len(result) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		result[i], result[j.Int64()] = result[j.Int64()], result[i]
	}

	return string(result), nil
}
