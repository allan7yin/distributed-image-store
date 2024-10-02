package main

import (
	"bit-image/internal/postrges"
	"bit-image/pkg/middleware"
	"bit-image/wire"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var handler *postrges.ConnectionHandler

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Postgres Initialization
	handler, err := postrges.NewConnectionHandler()
	if err != nil {
		log.Fatalf("Error initializing database connection: %v", err)
	}
	defer handler.Close()

	router := gin.Default()

	// healthcheck
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// db healthcheck
	router.GET("/db-stats", func(c *gin.Context) {
		handler.PrintConnectionPoolStats()
		c.JSON(http.StatusOK, gin.H{
			"message": "Database stats printed",
		})
	})

	// Initialize image handler
	imageHandler, err := wire.InitializeImageHandler()
	if err != nil {
		log.Fatalf("Failed to initialize the app: %v", err)
	}

	// Protected routes using AuthMiddleware
	apiGroup := router.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware())

	apiGroup.PUT("/generateUploadUrls", imageHandler.GeneratePresignedURL())
	apiGroup.POST("/confirmImageUploads", imageHandler.ConfirmImageUploads())

	// Start the server
	if err := router.Run(); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
