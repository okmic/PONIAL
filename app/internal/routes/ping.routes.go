package routes

import (
	"ponial/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupPingRoutes(router *gin.Engine, controller *controllers.PingController) {
	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/ping", controller.Ping)
		apiV1.Any("/ping/all", controller.PingAllMethods)
		apiV1.GET("/health", controller.HealthCheck)
	}
	router.GET("/", controller.Ping)
	router.Any("/ping", controller.PingAllMethods)
	router.Any("/all", controller.PingAllMethods)
}
