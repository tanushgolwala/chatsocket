package server

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientManager struct {
	clients map[string]*websocket.Conn
	lock    sync.Mutex
}

// NewClientManager initializes a new client manager
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[string]*websocket.Conn),
	}
}

// AddClient adds a new client to the manager
func (m *ClientManager) AddClient(clientID string, conn *websocket.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.clients[clientID] = conn
	fmt.Printf("Client %s connected\n", clientID)
}

// RemoveClient removes a client from the manager
func (m *ClientManager) RemoveClient(clientID string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.clients, clientID)
	fmt.Printf("Client %s disconnected\n", clientID)
}

// GetClient returns a client connection by ID
func (m *ClientManager) GetClient(clientID string) (*websocket.Conn, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	conn, exists := m.clients[clientID]
	return conn, exists
}
