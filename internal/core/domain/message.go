package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	DeletedAt        gorm.DeletedAt

	Chat           Chat
	User           User
	ReplyToMessage *Message
	Replies        []Message
}

type MessageRead struct {
	MessageID uuid.UUID
	UserID    uuid.UUID
	ReadAt    time.Time

	Message Message
	User    User
}

func (MessageRead) TableName() string {
	return "message_reads"
}
