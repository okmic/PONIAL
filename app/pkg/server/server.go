package server

import (
	"log"
	"time"

	"ponial/internal/controllers"
	"ponial/internal/database"
	"ponial/internal/repositories"
	"ponial/internal/routes"
	"ponial/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Host string
	Port string
	Mode string
}

type Server struct {
	config          *Config
	router          *gin.Engine
	pingController  *controllers.PingController
	usersController *controllers.UsersController
}

func New(config *Config) *Server {
	if config.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db := database.GetDB()

	if db == nil {
		log.Fatal("Database connection is not initialized. Make sure database.Connect() was called before server.New()")
	}
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	usersController := controllers.NewUsersController(userService)

	return &Server{
		config:          config,
		router:          router,
		pingController:  controllers.NewPingController(),
		usersController: usersController,
	}
}

func (s *Server) SetupRoutes() {
	routes.SetupPingRoutes(s.router, s.pingController)
	routes.SetupUsersRoutes(s.router, s.usersController)
	s.router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "Route not found",
			"path":    c.Request.URL.Path,
			"method":  c.Request.Method,
			"message": "Try /ping or /api/v1/users",
		})
	})
}

func (s *Server) Start() error {
	s.SetupRoutes()
	addr := s.config.Host + ":" + s.config.Port
	log.Printf("Server starting on %s", addr)
	return s.router.Run(addr)
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
