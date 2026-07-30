package utils

import (
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	cost := bcrypt.DefaultCost
	if costStr := os.Getenv("BCRYPT_COST"); costStr != "" {
		if parsedCost, err := strconv.Atoi(costStr); err == nil {
			cost = parsedCost
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hashedPassword), err
}

func ComparePassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
