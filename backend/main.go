package main

import (
	"encoding/json"
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

type Message struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Content string `json:"content"`
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

func (s *WebSocketServer) handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Read the unique ID of the client from query parameters
	clientID := r.URL.Query().Get("id")
	if clientID == "" {
		fmt.Println("Client ID not provided, closing connection")
		return
	}

	// Add the client to the map
	s.lock.Lock()
	s.clients[clientID] = conn
	s.lock.Unlock()
	fmt.Printf("Client %s connected\n", clientID)

	// Handle incoming messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Read error for client %s: %v\n", clientID, err)
			break
		}

		// Parse the incoming message
		var message Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			fmt.Println("Invalid message format:", err)
			continue
		}

		// Send the message to the target client
		s.sendMessageToClient(message.To, msg)
	}

	s.lock.Lock()
	delete(s.clients, clientID)
	s.lock.Unlock()
	fmt.Printf("Client %s disconnected\n", clientID)
}

// sendMessageToClient sends a message to a specific client
func (s *WebSocketServer) sendMessageToClient(clientID string, msg []byte) {
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

func main() {
	server := NewWebSocketServer()

	http.HandleFunc("/ws", server.handleConnections)
	fmt.Println("WebSocket server started on ws://localhost:8080/ws?id=<clientID>")
	err := http.ListenAndServe("0.0.0.0:8080", nil)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
