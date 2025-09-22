package ws

import (
	"log"
	"net/http"

	"github.com/HellEaglee/Golang-Chat/internal/core/port"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WS struct {
	hub *Hub
}

func NewWS(cs port.ChatService, ms port.MessageService) *WS {
	hub := NewHub(cs, ms)
	go hub.Run()
	return &WS{
		hub: hub,
	}
}

func (h *WS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	username := r.Context().Value("username").(string)

	if userID == "" {
		http.Error(w, "Unauthorized: user not authenticated", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		Hub:      h.hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		Username: username,
		Rooms:    make(map[string]bool),
	}

	h.hub.Register <- client
	go client.WritePump()
	go client.ReadPump()
}
