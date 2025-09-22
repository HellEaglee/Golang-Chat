package ws

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type RoomMessage struct {
	RoomID string
	Data   []byte
}

// Sent to client on connection
type ChatListResponse struct {
	Chats []ChatListItem `json:"chats"`
}

type ChatListItem struct {
	ID             uuid.UUID   `json:"id"`
	Name           *string     `json:"name"`
	IsGroup        bool        `json:"is_group"`
	LastMessage    string      `json:"last_message"`
	LastMessageAt  time.Time   `json:"last_message_at"`
	UnreadCount    int         `json:"unread_count,omitempty"`
	ParticipantIDs []uuid.UUID `json:"participant_ids,omitempty"`
}

// Sent when user joins chat
type MessageHistoryResponse struct {
	Messages []MessageItem `json:"messages"`
}

type MessageItem struct {
	ID               uuid.UUID     `json:"id"`
	ChatID           uuid.UUID     `json:"chat_id"`
	UserID           uuid.UUID     `json:"user_id"`
	Username         string        `json:"username"`
	Text             string        `json:"text"`
	IsEdited         bool          `json:"is_edited"`
	ReplyToMessageID *uuid.UUID    `json:"reply_to_message_id,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
	Replies          []MessageItem `json:"replies,omitempty"`
}

// Incoming from client
type ClientMessage struct {
	RoomID string `json:"room_id"`
	Text   string `json:"text"`
	// Optional: ReplyToMessageID *uuid.UUID
}
