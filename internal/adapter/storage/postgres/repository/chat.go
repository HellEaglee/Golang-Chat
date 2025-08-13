package repository

import (
	"context"
	"time"

	"github.com/HellEaglee/Golang-Chat/internal/adapter/storage/postgres"
	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
)

type ChatRepository struct {
	db *postgres.DB
}

func NewChatRepository(db *postgres.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// ----------------------------------------------------CHATS----------------------------------------------------
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

func (r *ChatRepository) GetChatsByUserID(ctx context.Context, id string) ([]domain.Chat, error) {
	var chats []domain.Chat
	if err := r.db.WithContext(ctx).Select("chats.*").Joins("JOIN chat_participants cp ON chats.id = cp.chat_id").Where("cp.user_id = ? AND chats.deleted_at IS NULL AND cp.deleted_at IS NULL", id).Order("chats.last_message_at DESC").Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
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

// ----------------------------------------------------CHAT_PARTICIPANTS----------------------------------------------------
func (r *ChatRepository) CreateChatParticipant(ctx context.Context, chatParticipant *domain.ChatParticipant) (*domain.ChatParticipant, error) {
	if err := r.db.WithContext(ctx).Create(chatParticipant).Error; err != nil {
		return nil, err
	}
	return chatParticipant, nil
}

func (r *ChatRepository) GetChatParticipantByChatIDUserID(ctx context.Context, chatID, userID string) (*domain.ChatParticipant, error) {
	var chatParticipant domain.ChatParticipant
	if err := r.db.WithContext(ctx).Where("chat_id = $1 AND user_id = $2 and deleted_at IS NULL", chatID, userID).First(&chatParticipant).Error; err != nil {
		return nil, err
	}
	return &chatParticipant, nil
}

func (r *ChatRepository) GetChatParticipantsByChatID(ctx context.Context, id string) ([]domain.ChatParticipant, error) {
	var chatParticipants []domain.ChatParticipant
	if err := r.db.WithContext(ctx).Where("chat_id = ?", id).Find(&chatParticipants).Error; err != nil {
		return nil, err
	}
	return chatParticipants, nil
}

func (r *ChatRepository) UpdateChatParticipant(ctx context.Context, chatParticipant *domain.ChatParticipant) (*domain.ChatParticipant, error) {
	var updatedChatParticipant domain.ChatParticipant
	query := `UPDATE chat_participants SET role = $3, updated_at = NOW() WHERE chat_id = $1 AND user_id = $2 AND deleted_at IS NULL RETURNING *`

	if err := r.db.WithContext(ctx).Raw(query, chatParticipant.ChatID, chatParticipant.UserID, chatParticipant.Role).Scan(&updatedChatParticipant).Error; err != nil {
		return nil, err
	}
	return &updatedChatParticipant, nil
}

func (r *ChatRepository) DeleteChatParticipant(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&domain.ChatParticipant{}).Where("id = ?", id).Updates(map[string]interface{}{
		"deleted_at": time.Now,
		"left_at":    time.Now,
	}).Error; err != nil {
		return err
	}
	return nil
}
