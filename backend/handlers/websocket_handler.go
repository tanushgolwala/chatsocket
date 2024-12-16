package handlers

import (
	"chatback/models"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Client struct {
	ID   uint
	Conn *websocket.Conn
	Send chan []byte
}

var clients = make(map[uint]*Client)
var clientsMutex sync.Mutex

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	// Example: Retrieve user ID from query params
	userID := uint(1) // Replace with actual logic

	client := &Client{ID: userID, Conn: conn, Send: make(chan []byte)}
	clientsMutex.Lock()
	clients[userID] = client
	clientsMutex.Unlock()

	go handleClientMessages(client, db)
}

func handleClientMessages(client *Client, db *gorm.DB) {
	defer func() {
		clientsMutex.Lock()
		delete(clients, client.ID)
		clientsMutex.Unlock()
		client.Conn.Close()
	}()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		var message struct {
			SenderID   uint   `json:"sender_id"`
			ReceiverID uint   `json:"receiver_id"`
			Content    string `json:"content"`
		}
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Failed to parse message:", err)
			continue
		}

		// Save message to DB (use GORM)
		saveMessage(db, message.SenderID, message.ReceiverID, message.Content)

		// Send message to the receiver if online
		clientsMutex.Lock()
		if receiver, ok := clients[message.ReceiverID]; ok {
			receiver.Send <- msg
		}
		clientsMutex.Unlock()
	}
}

func saveMessage(db *gorm.DB, senderID, receiverID uint, content string) {
	message := models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}
	db.Create(&message)
}
