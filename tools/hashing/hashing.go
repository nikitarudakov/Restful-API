package hashing

import (
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashUsername(username string) string {
	h := sha256.New()
	h.Write([]byte(username))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HashPassword(ac string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(ac), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func CheckPassword(hac string, ac string) error {
	return bcrypt.CompareHashAndPassword([]byte(hac), []byte(ac))
}
