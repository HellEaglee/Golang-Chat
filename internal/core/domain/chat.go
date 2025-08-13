package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chat struct {
	ID        uuid.UUID
	Name      *string
	IsGroup   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Participants []ChatParticipant
	Messages     []Message
}

type ChatParticipant struct {
	ChatID    uuid.UUID
	UserID    uuid.UUID
	Role      string
	JoinedAt  time.Time
	LeftAt    *time.Time
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Chat Chat
	User User
}

func (ChatParticipant) TableName() string {
	return "chat_participants"
}
