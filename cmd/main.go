package main

import (
	"bit-image/internal"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var handler *internal.ConnectionHandler

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Eagerly initialize the database connection when the app starts
	handler, err = internal.NewConnectionHandler()
	if err != nil {
		log.Fatalf("Error initializing database connection: %v", err)
	}
	defer handler.Close() // Ensure the connection pool is closed when the app shuts down

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
		// Print connection pool stats
		handler.PrintConnectionPoolStats()
		c.JSON(http.StatusOK, gin.H{
			"message": "Database stats printed",
		})
	})

	// Start the server
	if err := router.Run(); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
