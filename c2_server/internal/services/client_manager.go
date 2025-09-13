package services

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a connected client
type Client struct {
	ID   string
	Conn *websocket.Conn
}

// DefaultClientManager implements the ClientManager interface
type DefaultClientManager struct {
	clients    map[string]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.Mutex
}

// NewClientManager creates a new client manager and starts its processing loop
func NewClientManager() *DefaultClientManager {
	manager := &DefaultClientManager{
		clients:    make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	// Start the manager's processing loop
	go manager.run()

	return manager
}

// run processes client events in a loop
func (m *DefaultClientManager) run() {
	for {
		select {
		case client := <-m.register:
			m.mutex.Lock()
			m.clients[client.ID] = client
			m.mutex.Unlock()
			log.Printf("New client registered: %s", client.ID)

		case client := <-m.unregister:
			m.mutex.Lock()
			if _, ok := m.clients[client.ID]; ok {
				delete(m.clients, client.ID)
				log.Printf("Client unregistered: %s", client.ID)
			}
			m.mutex.Unlock()

		case message := <-m.broadcast:
			m.mutex.Lock()
			for _, client := range m.clients {
				if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("Error broadcasting to client %s: %v", client.ID, err)
				}
			}
			m.mutex.Unlock()
		}
	}
}

// RegisterClient registers a new client synchronously to ensure immediate availability.
func (m *DefaultClientManager) RegisterClient(client *Client) {
    m.mutex.Lock()
    m.clients[client.ID] = client
    m.mutex.Unlock()
    log.Printf("New client registered: %s", client.ID)
}

// UnregisterClient unregisters a client synchronously.
func (m *DefaultClientManager) UnregisterClient(client *Client) {
    m.mutex.Lock()
    if _, ok := m.clients[client.ID]; ok {
        delete(m.clients, client.ID)
        log.Printf("Client unregistered: %s", client.ID)
    }
    m.mutex.Unlock()
}

// BroadcastMessage sends a message to all connected clients
func (m *DefaultClientManager) BroadcastMessage(message []byte) {
	m.broadcast <- message
}

// SendToClient sends a message to a specific client
func (m *DefaultClientManager) SendToClient(clientID string, message []byte) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, ok := m.clients[clientID]
	if !ok {
		return fmt.Errorf("client not found: %s", clientID)
	}

	if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

// GetClients returns a map of all connected clients
func (m *DefaultClientManager) GetClients() map[string]*Client {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Create a copy of the clients map to avoid concurrent access issues
	clients := make(map[string]*Client, len(m.clients))
	for id, client := range m.clients {
		clients[id] = client
	}

	return clients
}