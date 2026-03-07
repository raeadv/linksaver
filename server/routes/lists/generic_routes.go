package routes

import (
	"linksaver/server/routes/handler"
	"linksaver/server/routes/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterGenericRoutes(u *gin.RouterGroup) {

	securedRoute := u.Use(middlewares.MustHaveValidToken())

	securedRoute.POST("/tags", handler.HandleCreateTag)
	securedRoute.GET("/tags", handler.HandleGetTags)
	securedRoute.GET("/tags/:query", handler.HandleSearchTags)
	securedRoute.DELETE("/tags/:id", handler.HandleDeleteTags)

	securedRoute.POST("/users", func(gc *gin.Context) {})

	securedRoute.GET("/links", handler.HandleGetLinks)
	securedRoute.GET("/links/get/web-meta", handler.HandleGetWebsiteMeta)
	securedRoute.POST("/links", handler.HandleCreateLink)

}
