package services

import (
	"log"

	"github.com/gorilla/websocket"
)

func NewClient(id string, conn *websocket.Conn) *Client {
	return &Client{
		ID:   id,
		Conn: conn,
	}
}

func (c *Client) HandleMessages() {
	defer func() {
		log.Printf("Client disconnected: %s", c.ID)
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("Read error from client %s: %v", c.ID, err)
			break
		}

		log.Printf("Received from %s: Type: %d, Message: %s", c.ID, messageType, string(p))
	}
}
