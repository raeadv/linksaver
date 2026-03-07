package utils

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var tokenLifetimeHrs time.Duration = time.Hour * 6

func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Dev fallback only. In production, always set JWT_SECRET.
		secret = "dev-only-change-this-secret-at-least-32-bytes"
	}
	return []byte(secret)

}

func GenerateToken(ID string, username string, email string) (string, error) {
	secretKey := getSecretKey()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":         ID,
			"username":   username,
			"email":      email,
			"expires_at": time.Now().Add(tokenLifetimeHrs).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func HashString(content string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(content), bcrypt.DefaultCost)
	return string(hashed), err
}

func ValidatePassword(expected string, provided string) error {
	err := bcrypt.CompareHashAndPassword([]byte(expected), []byte(provided))
	return err
}

func GetAuthHeader(gc *gin.Context) string {
	strAuthHeader := gc.GetHeader("Authorization")
	strToken := strings.TrimPrefix(strAuthHeader, "Bearer ")

	return strToken
}

func ValidateToken(strToken string) (*jwt.Token, error) {
	secretKey := getSecretKey()
	token, err := jwt.Parse(strToken, func(t *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid Token")
	}

	return token, nil
}
