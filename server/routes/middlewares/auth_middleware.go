package middlewares

import (
	"linksaver/server/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func MustHaveValidToken() gin.HandlerFunc {
	return func(gc *gin.Context) {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "dev-only-change-this-secret-at-least-32-bytes"
		}

		strToken := utils.GetAuthHeader(gc)
		if strToken == "" {
			gc.JSON(http.StatusForbidden, gin.H{
				"status":  false,
				"message": "invalid token [1]",
			})
			return
		}

		token, err := utils.ValidateToken(strToken)

		if err != nil {
			gc.JSON(http.StatusForbidden, gin.H{
				"status":  false,
				"message": "failed to validate token [2]",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		id, ok := claims["id"].(string)
		if !ok {
			gc.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  false,
				"message": "invalid token id [3]",
			})
			return
		}

		gc.Set("ID", id)

		gc.Next()
	}
}
