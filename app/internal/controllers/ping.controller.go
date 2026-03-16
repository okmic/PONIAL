package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingController struct{}

func NewPingController() *PingController {
	return &PingController{}
}

func (c *PingController) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"method":  ctx.Request.Method,
		"path":    ctx.Request.URL.Path,
		"status":  "success",
	})
}

func (c *PingController) PingAllMethods(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong from all methods",
		"method":  ctx.Request.Method,
		"path":    ctx.Request.URL.Path,
		"status":  "success",
	})
}

func (c *PingController) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "ping-service",
		"version": "1.0.0",
	})
}
