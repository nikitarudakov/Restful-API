package hashing

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashAuthCredentials(ac string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(ac), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func CheckAuthCredentials(hac string, ac string) error {
	return bcrypt.CompareHashAndPassword([]byte(hac), []byte(ac))
}
