package services

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
)

func setupTestWebSocket(t *testing.T) (*websocket.Conn, *httptest.Server, error) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer conn.Close()

		// Echo back any messages received
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			if err := conn.WriteMessage(msgType, msg); err != nil {
				break
			}
		}
	}))

	// Convert http://127.0.0.1... to ws://127.0.0.1...
	wsURL := "ws" + server.URL[4:]

	// Connect to the test server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)

	return conn, server, err
}

func TestClientManager(t *testing.T) {
	// Create a new client manager
	manager := NewClientManager()

	// Test client registration
	t.Run("RegisterClient", func(t *testing.T) {
		conn, server, err := setupTestWebSocket(t)
		if err != nil {
			t.Fatalf("Failed to create test WebSocket: %v", err)
		}
		defer server.Close()

		client := &Client{ID: "test-client", Conn: conn}
		manager.RegisterClient(client)

		clients := manager.GetClients()
		if _, ok := clients["test-client"]; !ok {
			t.Errorf("Client was not registered")
		}
	})

	t.Run("UnregisterClient", func(t *testing.T) {
		conn, server, err := setupTestWebSocket(t)

		if err != nil {
			t.Fatalf("Failed to create test WebSocket: %v", err)
		}
		defer server.Close()

		client := &Client{ID: "test-client-2", Conn: conn}
		manager.RegisterClient(client)
		manager.UnregisterClient(client)

		clients := manager.GetClients()
		if _, ok := clients["test-client-2"]; ok {
			t.Errorf("Client was not unregistered")
		}
	})

	t.Run("SendToClient", func(t *testing.T) {
		conn, server, err := setupTestWebSocket(t)
		if err != nil {
			t.Fatalf("Failed to create test WebSocket: %v", err)
		}
		defer server.Close()

		client := &Client{ID: "test-client-3", Conn: conn}
		manager.RegisterClient(client)

		err = manager.SendToClient("test-client-3", []byte("test message"))
		if err != nil {
			t.Errorf("Failed to send message to client: %v", err)
		}

		// Send to a non-existent client should fail
		err = manager.SendToClient("non-existent", []byte("test message"))
		if err == nil {
			t.Errorf("Expected error when sending to non-existent client")
		}
	})
}
