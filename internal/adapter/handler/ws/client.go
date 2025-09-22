package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte

	UserID   string
	Username string

	mu    sync.RWMutex
	Rooms map[string]bool
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("%v", err)
			continue
		}

		c.Hub.HandleMessage(c, msg)
	}
}

func (c *Client) WritePump() {
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func (c *Client) JoinRoom(roomID string) {
	c.mu.Lock()
	c.Rooms[roomID] = true
	c.mu.Unlock()

	joinMsg, _ := json.Marshal(Message{
		Type: "system",
		Payload: map[string]string{
			"event":    "user_joined",
			"username": c.Username,
			"room_id":  roomID,
		},
	})

	c.Hub.BroadcastToRoom <- &RoomMessage{
		RoomID: roomID,
		Data:   joinMsg,
	}
}

func (c *Client) LeaveRoom(roomID string) {
	c.mu.Lock()
	delete(c.Rooms, roomID)
	c.mu.Unlock()

	leaveMsg, _ := json.Marshal(Message{
		Type: "system",
		Payload: map[string]string{
			"event":    "user_left",
			"username": c.Username,
			"room_id":  roomID,
		},
	})

	c.Hub.BroadcastToRoom <- &RoomMessage{
		RoomID: roomID,
		Data:   leaveMsg,
	}
}

func (c *Client) InRoom(roomID string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Rooms[roomID]
}
