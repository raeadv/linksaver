package routes

import (
	routes "linksaver/server/routes/lists"

	"github.com/gin-gonic/gin"
)

func RegisterApiRoute(api *gin.RouterGroup) {

	routes.RegisterAuthRoutes(api.Group("/auth"))

	routes.RegisterGenericRoutes(api)

}
