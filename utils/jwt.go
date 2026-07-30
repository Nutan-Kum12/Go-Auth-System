package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ACCESSSECRETKEY = []byte(os.Getenv("ACCESS_SECRET"))
var REFRESHSECRETKEY = []byte(os.Getenv("REFRESH_SECRET"))

func generateToken(email string, duration time.Duration, tokenType string, secret []byte) (string, error) {
	claims := jwt.MapClaims{ //create claims
		"email": email,
		"type":  tokenType,
		"exp":   time.Now().Add(duration).Unix(),
		"iat":   time.Now().Unix(),
		"sub":   email,
	}
	token := jwt.NewWithClaims( ////Put claims into JWT
		jwt.SigningMethodHS256,
		claims,
	)
	return token.SignedString(secret) //sign the JWT
}

func GenerateAccessToken(email string) (string, error) {
	return generateToken(email, 15*time.Minute, "access", ACCESSSECRETKEY)
}

func GenerateRefreshToken(email string) (string, error) {
	return generateToken(email, 7*24*time.Hour, "refresh", REFRESHSECRETKEY)
}

// jwt validation(extract claims from jwt send by client(user))
// tokenString receive from the client
// Parse & Verify JWT
func ParseToken(tokenString string, secret []byte) (*jwt.Token, error) {
	return jwt.Parse(tokenString,
		func(token *jwt.Token) (any, error) {
			// Check Signing Method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Secret used for verification
			return secret, nil
		},
	)
}

// Wrapper for Access Token
func ParseAccessToken(tokenString string) (*jwt.Token, error) {
	return ParseToken(tokenString, ACCESSSECRETKEY)
}

// Wrapper for Refresh Token
func ParseRefreshToken(tokenString string) (*jwt.Token, error) {
	return ParseToken(tokenString, REFRESHSECRETKEY)
}
