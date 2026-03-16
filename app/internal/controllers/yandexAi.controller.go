package controllers

import (
	"net/http"
	"strconv"

	"ponial/internal/models"
	"ponial/internal/services"

	"github.com/gin-gonic/gin"
)

type YandexAIController struct {
	yandexAIService services.YandexAIService
}

func NewYandexAIController(yandexAIService services.YandexAIService) *YandexAIController {
	return &YandexAIController{yandexAIService: yandexAIService}
}

func (c *YandexAIController) CreateYandexAI(ctx *gin.Context) {
	var req models.YandexAICreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	yandexAI, err := c.yandexAIService.CreateYandexAI(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, yandexAI)
}

func (c *YandexAIController) GetYandexAIByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid yandex AI ID"})
		return
	}

	yandexAI, err := c.yandexAIService.GetYandexAI(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, yandexAI)
}

func (c *YandexAIController) GetYandexAIByToken(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	yandexAI, err := c.yandexAIService.GetYandexAIByToken(token)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, yandexAI)
}

func (c *YandexAIController) UpdateYandexAI(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid yandex AI ID"})
		return
	}

	var req models.YandexAIUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	yandexAI, err := c.yandexAIService.UpdateYandexAI(id, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, yandexAI)
}

func (c *YandexAIController) DeleteYandexAI(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid yandex AI ID"})
		return
	}

	if err := c.yandexAIService.DeleteYandexAI(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *YandexAIController) GetAllYandexAIs(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	aiType := ctx.Query("type")
	token := ctx.Query("token")
	folderID := ctx.Query("folderId")
	minPrice, _ := strconv.ParseInt(ctx.Query("minPrice"), 10, 64)
	maxPrice, _ := strconv.ParseInt(ctx.Query("maxPrice"), 10, 64)

	yandexAIs, err := c.yandexAIService.ListYandexAIs(limit, offset, aiType, token, folderID, minPrice, maxPrice)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, yandexAIs)
}

func (c *YandexAIController) GetYandexAIsByType(ctx *gin.Context) {
	aiType := ctx.Param("type")
	if aiType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "type is required"})
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	token := ctx.Query("token")
	folderID := ctx.Query("folderId")
	minPrice, _ := strconv.ParseInt(ctx.Query("minPrice"), 10, 64)
	maxPrice, _ := strconv.ParseInt(ctx.Query("maxPrice"), 10, 64)

	yandexAIs, err := c.yandexAIService.ListYandexAIs(limit, offset, aiType, token, folderID, minPrice, maxPrice)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, yandexAIs)
}
