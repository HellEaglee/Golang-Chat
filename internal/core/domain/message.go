package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID               uuid.UUID
	ChatID           uuid.UUID
	UserID           uuid.UUID
	Text             string
	IsEdited         bool
	ReplyToMessageID *uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time

	Chat           Chat
	User           User
	ReplyToMessage *Message
	Replies        []Message
}
