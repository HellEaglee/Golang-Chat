package port

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
)

type ChatRepository interface {
	CreateChat(ctx context.Context, chat *domain.Chat) (*domain.Chat, error)
	GetChatByID(ctx context.Context, id string) (*domain.Chat, error)
	GetChats(ctx context.Context, skip uint64, limit uint64) ([]domain.Chat, error)
	UpdateChat(ctx context.Context, chat *domain.Chat) (*domain.Chat, error)
	DeleteChat(ctx context.Context, id string) error
}
