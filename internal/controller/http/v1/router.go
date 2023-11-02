package v1

import (
	_ "usersegments/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
	)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		h.UserHandler.newUserRoutes(api)
		h.SegmentHandler.newSegmentRoutes(api)
		h.OperationHandler.newOperationRoutes(api)

	}
}
