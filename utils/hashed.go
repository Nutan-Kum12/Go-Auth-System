package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		4, // we can use more like 8, 10, 12 and it takes more CPU usage and high security , 4 gives less security and use less CPU
	)
	return string(hashedPassword), err
}
func ComparePassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
