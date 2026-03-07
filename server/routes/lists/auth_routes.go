package routes

import (
	"linksaver/server/routes/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(auth_r *gin.RouterGroup) {

	auth_r.POST("/login", handler.HandleLogin)
	auth_r.POST("/register", handler.HandleRegister)
	auth_r.GET("/refresh-tokenn", handler.HandleRefreshToken)

}
