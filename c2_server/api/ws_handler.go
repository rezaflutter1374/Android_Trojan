package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"hermes_spectre/c2_server/internal/services"
)

// WebSocket upgrader configuration
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for demonstration purposes
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WSHandler handles WebSocket connections from clients
func WSHandler(clientManager services.ClientManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client ID from query parameters
		clientID := c.Query("id")
		if clientID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Client ID is required"})
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection: %v", err)
			return
		}

		// Register the client with the client manager
		client := services.NewClient(clientID, conn)
		clientManager.RegisterClient(client)

		// Start handling messages in a goroutine
		go client.HandleMessages()

		log.Printf("Client connected: %s", clientID)
	}
}