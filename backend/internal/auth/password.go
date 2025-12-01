package auth

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	// MinPasswordLength minimum password length
	MinPasswordLength = 12
	// BcryptCost for password hashing
	BcryptCost = 12
	// MaxFailedLogins before account lockout
	MaxFailedLogins = 5
)

var (
	ErrWeakPassword    = errors.New("password does not meet complexity requirements")
	ErrInvalidPassword = errors.New("invalid password")
	ErrAccountLocked   = errors.New("account is locked due to too many failed login attempts")
)

// HashPassword hashes a plaintext password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword checks if a password matches its hash
func VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordPolicy enforces password complexity requirements
// Requirements: min 12 chars, 1 uppercase, 1 number, 1 symbol
func ValidatePasswordPolicy(password string) error {
	if len(password) < MinPasswordLength {
		return ErrWeakPassword
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return ErrWeakPassword
	}

	return nil
}
