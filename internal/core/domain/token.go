package domain

import (
	"time"

	"github.com/google/uuid"
)

type TokenPayload struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

type Token struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time
}
