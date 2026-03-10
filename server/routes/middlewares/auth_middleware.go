package middlewares

import (
	"linksaver/server/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

		_, err := utils.ValidateToken(strToken, false)
		if err != nil {
			gc.JSON(http.StatusForbidden, gin.H{
				"status":  false,
				"message": "failed to validate token [2]",
			})
			return
		}

		gc.Next()
	}
}
