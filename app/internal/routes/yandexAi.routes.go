package routes

import (
	"ponial/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupYandexAIRoutes(router *gin.Engine, controller *controllers.YandexAIController) {
	yandexAI := router.Group("/api/v1/yandex-ai")
	{
		yandexAI.GET("", controller.GetAllYandexAIs)
		yandexAI.POST("", controller.CreateYandexAI)
		yandexAI.GET("/type/:type", controller.GetYandexAIsByType)
		yandexAI.GET("/token/:token", controller.GetYandexAIByToken)

		ai := yandexAI.Group("/:id")
		{
			ai.GET("", controller.GetYandexAIByID)
			ai.PUT("", controller.UpdateYandexAI)
			ai.PATCH("", controller.UpdateYandexAI)
			ai.DELETE("", controller.DeleteYandexAI)
		}
	}
}
