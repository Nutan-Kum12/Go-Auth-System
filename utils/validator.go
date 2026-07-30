package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func ValidateEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	//Standard regex for email validation
	return regexp.MustCompile(regex).MatchString(email)
}
func ValidatePassword(password string,) bool {

	// minimum length
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	specialChars := "!@#$%^&*"

	for _, ch := range password {

		if unicode.IsUpper(ch) {
			hasUpper = true
		}

		if unicode.IsLower(ch) {
			hasLower = true
		}

		if unicode.IsDigit(ch) {
			hasNumber = true
		}

		if strings.ContainsRune(
			specialChars,
			ch,
		) {
			hasSpecial = true
		}
	}

	return hasUpper &&
		hasLower &&
		hasNumber &&
		hasSpecial
}
func ValidateName(name string) bool {
	regex := `^[a-zA-Z ]{2,40}$`
	//at least 2 characters,max 40 characters,only letters and spaces
	return regexp.MustCompile(regex).MatchString(name)
}
