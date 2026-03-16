package routes

import (
	"ponial/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupWorkspaceRoutes(router *gin.Engine, controller *controllers.WorkspaceController) {
	workspaces := router.Group("/api/v1/workspaces")
	{
		workspaces.GET("/me", controller.GetMyWorkspace)
		workspaces.PUT("/me", controller.UpdateMyWorkspace)
		workspaces.GET("", controller.ListWorkspaces)
		workspace := workspaces.Group("/:id")
		{
			workspace.GET("", controller.GetWorkspace)
			workspace.PUT("", controller.UpdateWorkspace)
		}
	}
}
