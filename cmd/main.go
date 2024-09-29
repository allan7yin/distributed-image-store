package main

import (
	"bit-image/internal/postrges"
	"bit-image/pkg/handlers"
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
	handler, err = postrges.NewConnectionHandler()
	if err != nil {
		log.Fatalf("Error initializing database connection: %v", err)
	}
	defer handler.Close()

	// Create a new Gin router
	router := gin.Default()

	// Define routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Example route that uses the database connection
	router.GET("/db-stats", func(c *gin.Context) {
		handler.PrintConnectionPoolStats()
		c.JSON(http.StatusOK, gin.H{
			"message": "Database stats printed",
		})
	})

	imageService, err := wire.WireApp()
	if err != nil {
		log.Fatalf("Failed to initialize the app: %v", err)
	}

	imageHandler := handlers.NewImageHandler(imageService)
	router.PUT("/generateUploadUrls", imageHandler.GeneratePresignedURL())
	router.POST("/confirmImageUploads", imageHandler.ConfirmImageUploads())

	// Start the server
	if err := router.Run(); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
