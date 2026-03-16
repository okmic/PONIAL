package main

import (
	"log"
	"ponial/internal/database"
	"ponial/pkg/config"
	"ponial/pkg/server"
)

func main() {
	cfg := config.MustLoad()
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	println("dsad")
	defer database.Close()
	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database initialized successfully")
	app := &server.Config{
		Host: cfg.AppHost,
		Port: cfg.AppPort,
		Mode: cfg.AppMode,
	}

	srv := server.New(app)
	if err := srv.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
