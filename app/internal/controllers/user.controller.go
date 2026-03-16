package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"ponial/internal/models"
	"ponial/internal/services"
	"ponial/pkg/config"
	"ponial/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	userService services.UserService
}

func NewUsersController(userService services.UserService) *UsersController {
	return &UsersController{userService: userService}
}

func (c *UsersController) Signup(ctx *gin.Context) {
	adminSecret := config.MustLoad().AdminSecret
	var req models.UserSignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.AdminSecret != adminSecret {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Property AdminSecret is not valid"})
		return
	}
	response, err := c.userService.Signup(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *UsersController) Signin(ctx *gin.Context) {
	var req models.UserSigninRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.userService.Signin(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *UsersController) GetAllUsers(ctx *gin.Context) {
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

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	role := ctx.Query("role")
	name := ctx.Query("name")

	users, err := c.userService.ListUsers(limit, offset, role, name, models.Role(claims.Role), claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (c *UsersController) CreateUser(ctx *gin.Context) {
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

	creator, err := c.userService.GetUser(claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "creator not found"})
		return
	}

	var req models.UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.CreateUser(&req, models.Role(creator.Role), creator.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c *UsersController) GetUserByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := c.userService.GetUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UsersController) GetMe(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "authorization header is required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		tokenString = authHeader
	}

	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	user, err := c.userService.GetUser(claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UsersController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.UpdateUser(id, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UsersController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := c.userService.DeleteUser(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *UsersController) GetUsersByRole(ctx *gin.Context) {
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

	role := ctx.Param("role")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	name := ctx.Query("name")

	users, err := c.userService.ListUsers(limit, offset, role, name, models.Role(claims.Role), claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}
