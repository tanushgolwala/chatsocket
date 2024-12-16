package routes

import (
	"chatback/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupWebSocketRoutes(router *mux.Router, db *gorm.DB) {
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.WebSocketHandler(w, r, db)
	})
}
