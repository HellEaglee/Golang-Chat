package domain

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID
	Name      *string
	IsGroup   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	Participants []ChatParticipant
	Messages     []Message
}

type ChatParticipant struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	UserID    uuid.UUID
	Role      string
	JoinedAt  time.Time
	LeftAt    *time.Time
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time

	Chat Chat
	User User
}
