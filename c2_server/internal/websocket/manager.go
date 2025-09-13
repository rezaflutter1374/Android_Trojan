package websocket

import (
	"log"
	"net"
	"sync"
)

type ClientManager struct {
	clients    map[string]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.Mutex
}

func NewManager() *ClientManager {
	return &ClientManager{
		clients:    make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (m *ClientManager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mutex.Lock()
			m.clients[client.ID] = client
			m.mutex.Unlock()
			log.Printf("INFO: New client registered: %s (Total clients: %d)", client.ID, len(m.clients))

		case client := <-m.unregister:
			m.mutex.Lock()
			if _, ok := m.clients[client.ID]; ok {

				if conn := m.clients[client.ID].Conn; conn != nil {
					conn.Close()
				}
				delete(m.clients, client.ID)
				log.Printf("INFO: Client unregistered: %s (Total clients: %d)", client.ID, len(m.clients))
			}
			m.mutex.Unlock()

		case message := <-m.broadcast:
			m.mutex.Lock()

			for id, client := range m.clients {
				if err := client.Conn.WriteMessage(1, message); err != nil { // 1 for TextMessage
					log.Printf("ERROR: Failed to broadcast to client %s: %v", id, err)
				}
			}
			m.mutex.Unlock()
		}
	}
}

func (m *ClientManager) SendToClient(clientID string, message []byte) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, ok := m.clients[clientID]
	if !ok {
		return &net.AddrError{Err: "Client not found", Addr: clientID}
	}

	return client.Conn.WriteMessage(1, message) // 1 for TextMessage
}

func (m *ClientManager) GetRegisterChan() chan *Client {
	return m.register
}

// GetUnregisterChan کانال unregister را برای استفاده خارجی برمی‌گرداند.
func (m *ClientManager) GetUnregisterChan() chan *Client {
	return m.unregister
}

func (m *ClientManager) BroadcastToAll(message []byte) {
	m.broadcast <- message
}
