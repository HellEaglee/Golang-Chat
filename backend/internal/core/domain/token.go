package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPayload struct {
	ID        uuid.UUID
	UserID    string
	SessionID string
}

type TokenClaims struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

type Token struct {
	ID        uuid.UUID
	UserID    string
	SessionID string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time
}
