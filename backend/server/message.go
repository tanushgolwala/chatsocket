package server

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Content string `json:"content"`
}

// ParseMessage parses a JSON-encoded message
func ParseMessage(data []byte) (*Message, error) {
	var message Message
	err := json.Unmarshal(data, &message)
	if err != nil {
		return nil, fmt.Errorf("invalid message format: %w", err)
	}
	return &message, nil
}
