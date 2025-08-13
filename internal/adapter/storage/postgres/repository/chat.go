package repository

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/adapter/storage/postgres"
	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
)

type ChatRepository struct {
	db *postgres.DB
}

func NewChatRepository(db *postgres.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateChat(ctx context.Context, chat *domain.Chat) (*domain.Chat, error) {
	if err := r.db.WithContext(ctx).Create(chat).Error; err != nil {
		return nil, err
	}
	return chat, nil
}

func (r *ChatRepository) GetChatByID(ctx context.Context, id string) (*domain.Chat, error) {
	var chat domain.Chat
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&chat).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *ChatRepository) GetChats(ctx context.Context, skip uint64, limit uint64) ([]domain.Chat, error) {
	var chats []domain.Chat
	if err := r.db.WithContext(ctx).Limit(int(limit)).Offset(int(skip)).Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *ChatRepository) UpdateChat(ctx context.Context, chat *domain.Chat) (*domain.Chat, error) {
	var updatedChat domain.Chat
	query := `UPDATE chats SET name = $2, is_group = $3, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL RETURNING *`

	if err := r.db.WithContext(ctx).Raw(query, chat.ID, chat.Name, chat.IsGroup).Scan(&updatedChat).Error; err != nil {
		return nil, err
	}
	return &updatedChat, nil
}

func (r *ChatRepository) DeleteChat(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Chat{}).Error; err != nil {
		return err
	}
	return nil
}
