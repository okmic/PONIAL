package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"ponial/internal/models"
	"ponial/internal/services"
	"ponial/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type WorkspaceController struct {
	workspaceService services.WorkspaceService
}

func NewWorkspaceController(workspaceService services.WorkspaceService) *WorkspaceController {
	return &WorkspaceController{workspaceService: workspaceService}
}

func (c *WorkspaceController) GetMyWorkspace(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	workspace, err := c.workspaceService.GetMyWorkspace(claims.UserID, models.Role(claims.Role))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, workspace)
}

func (c *WorkspaceController) UpdateMyWorkspace(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	var req models.WorkspaceUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workspace, err := c.workspaceService.UpdateMyWorkspace(&req, claims.UserID, models.Role(claims.Role))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, workspace)
}

func (c *WorkspaceController) GetWorkspace(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace ID"})
		return
	}

	workspace, err := c.workspaceService.GetWorkspace(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, workspace)
}

func (c *WorkspaceController) UpdateWorkspace(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace ID"})
		return
	}

	var req models.WorkspaceUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workspace, err := c.workspaceService.UpdateWorkspace(id, &req, claims.UserID, models.Role(claims.Role))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, workspace)
}

func (c *WorkspaceController) ListWorkspaces(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	name := ctx.Query("name")
	userHeadID, _ := strconv.ParseInt(ctx.Query("user_head_id"), 10, 64)

	workspaces, err := c.workspaceService.ListWorkspaces(limit, offset, name, userHeadID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, workspaces)
}
