package services

// ClientManager defines the interface for managing WebSocket clients
type ClientManager interface {
	// RegisterClient registers a new client connection
	RegisterClient(client *Client)
	
	// UnregisterClient removes a client connection
	UnregisterClient(client *Client)
	
	// BroadcastMessage sends a message to all connected clients
	BroadcastMessage(message []byte)
	
	// SendToClient sends a message to a specific client
	SendToClient(clientID string, message []byte) error
	
	// GetClients returns a copy of the map of all connected clients
	GetClients() map[string]*Client
	
	// Run starts the client manager's event loop
	// Note: DefaultClientManager implements this as a private 'run' method that is started automatically by NewClientManager
}