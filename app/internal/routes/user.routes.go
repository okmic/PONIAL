package routes

import (
	"ponial/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupUsersRoutes(router *gin.Engine, controller *controllers.UsersController) {
	router.POST("/api/v1/auth/signup", controller.Signup)
	router.POST("/api/v1/auth/signin", controller.Signin)

	users := router.Group("/api/v1/users")
	{
		users.GET("", controller.GetAllUsers)
		users.GET("/me", controller.GetMe)
		users.POST("", controller.CreateUser)
		users.GET("/role/:role", controller.GetUsersByRole)
		user := users.Group("/:id")
		{
			user.GET("", controller.GetUserByID)
			user.PUT("", controller.UpdateUser)
			user.PATCH("", controller.UpdateUser)
			user.DELETE("", controller.DeleteUser)
		}
	}
}
