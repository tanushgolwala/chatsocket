package server

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func SendMessageToClient(clientID string, msg []byte, clients map[string]*websocket.Conn) {
	conn, exists := clients[clientID]
	if !exists {
		fmt.Printf("Client %s not found\n", clientID)
		return
	}

	err := conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		fmt.Printf("Write error for client %s: %v\n", clientID, err)
		conn.Close()
		delete(clients, clientID)
	}
}
