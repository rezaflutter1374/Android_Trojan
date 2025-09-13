package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// Command represents a command received from the C2 server
type Command struct {
	Action  string `json:"action"`
	Payload string `json:"payload"`
}

// Response represents a response sent back to the C2 server
type Response struct {
	ClientID string      `json:"client_id"`
	Action   string      `json:"action"`
	Success  bool        `json:"success"`
	Result   interface{} `json:"result"`
	Error    string      `json:"error,omitempty"`
}

func main() {
	addr := flag.String("addr", "localhost:8080", "C2 server address")
	id := flag.String("id", "", "Client ID (optional)")
	flag.Parse()

	// Generate a client ID if not provided
	clientID := *id
	if clientID == "" {
		clientID = fmt.Sprintf("client-%d", time.Now().Unix())
	}

	log.Printf("Connecting to C2 server as client: %s", clientID)

	// Create WebSocket URL with client ID
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws", RawQuery: fmt.Sprintf("id=%s", clientID)}
	log.Printf("Connecting to %s", u.String())

	// Connect to the WebSocket server
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket server:", err)
	}
	defer c.Close()

	// Set up channel to handle OS signals for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Set up channel for outgoing messages
	done := make(chan struct{})

	// Start goroutine to read messages from the server
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}

			// Parse the command
			var cmd Command
			if unmarshalErr := json.Unmarshal(message, &cmd); unmarshalErr != nil {
				log.Printf("Failed to parse command: %v", err)
				continue
			}

			log.Printf("Received command: %s with payload: %s", cmd.Action, cmd.Payload)

			// Process the command (simulated)
			response := Response{
				ClientID: clientID,
				Action:   cmd.Action,
				Success:  true,
				Result:   fmt.Sprintf("Executed %s with %s", cmd.Action, cmd.Payload),
			}

			// Send response back to server
			respJSON, err := json.Marshal(response)
			if err != nil {
				log.Printf("Failed to marshal response: %v", err)
				continue
			}

			if err := c.WriteMessage(websocket.TextMessage, respJSON); err != nil {
				log.Println("Write error:", err)
				return
			}

			log.Printf("Sent response for command: %s", cmd.Action)
		}
	}()

	// Send periodic heartbeat
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// Send heartbeat
			heartbeat := Response{
				ClientID: clientID,
				Action:   "HEARTBEAT",
				Success:  true,
				Result:   map[string]interface{}{"status": "alive", "timestamp": time.Now().Unix()},
			}

			heartbeatJSON, err := json.Marshal(heartbeat)
			if err != nil {
				log.Printf("Failed to marshal heartbeat: %v", err)
				continue
			}

			if err := c.WriteMessage(websocket.TextMessage, heartbeatJSON); err != nil {
				log.Println("Write error:", err)
				return
			}

			log.Println("Sent heartbeat")
		case <-interrupt:
			log.Println("Received interrupt signal, closing connection")

			// Send close message to server
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Write close error:", err)
				return
			}

			// Wait for server to close the connection
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
