package main

import (
	"chatback/server"
	"fmt"
	"net/http"
)

func main() {
	server := server.NewWebSocketServer()

	http.HandleFunc("/ws", server.HandleConnections)
	fmt.Println("WebSocket server started on ws://localhost:8080/ws?id=<clientID>")
	err := http.ListenAndServe("0.0.0.0:8080", nil)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
