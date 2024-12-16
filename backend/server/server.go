package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketServer struct {
	clients map[string]*websocket.Conn
	lock    sync.Mutex
}

// NewWebSocketServer initializes the server
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients: make(map[string]*websocket.Conn),
	}
}

// HandleConnections handles the WebSocket connections
func (s *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	clientID := r.URL.Query().Get("id")
	if clientID == "" {
		fmt.Println("Client ID not provided, closing connection")
		return
	}

	s.lock.Lock()
	s.clients[clientID] = conn
	s.lock.Unlock()
	fmt.Printf("Client %s connected\n", clientID)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Read error for client %s: %v\n", clientID, err)
			break
		}

		// Parse the incoming message
		message, err := ParseMessage(msg)
		if err != nil {
			fmt.Println("Invalid message format:", err)
			continue
		}

		// Send the message to the target client
		s.SendMessageToClient(message.To, msg)
	}

	s.lock.Lock()
	delete(s.clients, clientID)
	s.lock.Unlock()
	fmt.Printf("Client %s disconnected\n", clientID)
}

// SendMessageToClient sends a message to a specific client
func (s *WebSocketServer) SendMessageToClient(clientID string, msg []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()

	conn, exists := s.clients[clientID]
	if !exists {
		fmt.Printf("Client %s not found\n", clientID)
		return
	}

	err := conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		fmt.Printf("Write error for client %s: %v\n", clientID, err)
		conn.Close()
		delete(s.clients, clientID)
	}
}
