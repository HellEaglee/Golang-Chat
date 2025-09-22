package port

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
)

type MessageRepository interface {
	// Message
	CreateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error)
	GetMessageByID(ctx context.Context, id string) (*domain.Message, error)
	GetMessagesByChatID(ctx context.Context, chatID string) ([]domain.Message, error)
	UpdateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error)
	DeleteMessage(ctx context.Context, id string) error
	// MessageRead
	CreateMessageRead(ctx context.Context, messageRead *domain.MessageRead) (*domain.MessageRead, error)
	GetMessageReadsByMessageID(ctx context.Context, id string) ([]domain.MessageRead, error)
	DeleteMessageRead(ctx context.Context, messageID, userID string) error
}

type MessageService interface {
	// Message
	CreateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error)
	GetMessage(ctx context.Context, id string) (*domain.Message, error)
	GetMessagesByChatID(ctx context.Context, chatID string) ([]domain.Message, error)
	UpdateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error)
	DeleteMessage(ctx context.Context, id string) error
	// MessageRead
	CreateMessageRead(ctx context.Context, messageRead *domain.MessageRead) (*domain.MessageRead, error)
	GetMessageReadsByMessageID(ctx context.Context, id string) ([]domain.MessageRead, error)
	DeleteMessageRead(ctx context.Context, messageID, userID string) error
}
