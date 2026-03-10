package utils

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var tokenLifetimeHrs time.Duration = time.Hour * 12

func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Dev fallback only. In production, always set JWT_SECRET.
		secret = "dev-only-change-this-secret-at-least-32-bytes"
	}
	return []byte(secret)

}

func GenerateToken(gc *gin.Context, ID string, username string, email string) (string, error) {
	secretKey := getSecretKey()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"email":    email,
			"exp":      time.Now().Add(tokenLifetimeHrs).Unix(),
		},
	)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       ID,
			"username": username,
			"email":    email,
			"exp":      time.Now().Add(tokenLifetimeHrs).Add(1 * time.Hour).Unix(),
		},
	)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	cookieExpire := time.Now().Add(tokenLifetimeHrs).Add(1 * time.Hour)

	gc.SetCookieData(
		&http.Cookie{
			Name:     "session_id",
			Value:    ID,
			Path:     "/",
			Expires:  cookieExpire,
			HttpOnly: true,
			Secure:   os.Getenv("MODE") == "production",
			SameSite: http.SameSiteLaxMode,
		},
	)

	gc.SetCookieData(
		&http.Cookie{
			Name:     "refresh_token",
			Value:    refreshTokenString,
			Path:     "/",
			Expires:  cookieExpire,
			HttpOnly: true,
			Secure:   os.Getenv("MODE") == "production",
			SameSite: http.SameSiteLaxMode,
		},
	)

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

func GetAuthSessionData(gc *gin.Context) (string, error) {
	refreshToken, err := gc.Cookie("refresh_token")
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func ValidateToken(strToken string, withoutValidation bool) (*jwt.Token, error) {
	secretKey := getSecretKey()
	if withoutValidation {
		token, err := jwt.Parse(strToken, func(t *jwt.Token) (any, error) {
			return []byte(secretKey), nil
		}, jwt.WithoutClaimsValidation())
		if err != nil {
			return nil, err
		}

		if !token.Valid {
			return nil, errors.New("Invalid Token")
		}

		return token, nil

	} else {
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

}
