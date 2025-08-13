package repository

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/adapter/storage/postgres"
	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
)

type MessageRepository struct {
	db *postgres.DB
}

func NewMessageRepository(db *postgres.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// ----------------------------------------------------MESSAGES----------------------------------------------------
func (r *MessageRepository) CreateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error) {
	if err := r.db.WithContext(ctx).Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

func (r *MessageRepository) GetMessageByID(ctx context.Context, id string) (*domain.Message, error) {
	var message domain.Message
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&message).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepository) GetMessagesByChatID(ctx context.Context, chatID string) ([]domain.Message, error) {
	var messages []domain.Message
	if err := r.db.WithContext(ctx).Where("chat_id = ? AND deleted_at IS NULL", chatID).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageRepository) UpdateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error) {
	var updatedMessage domain.Message
	query := `UPDATE messages SET text = $2, is_edited = TRUE, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL RETURNING *`

	if err := r.db.WithContext(ctx).Raw(query, message.ChatID, message.Text).Scan(&updatedMessage).Error; err != nil {
		return nil, err
	}
	return &updatedMessage, nil
}

func (r *MessageRepository) DeleteMessage(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Message{}).Error; err != nil {
		return err
	}
	return nil
}

// ----------------------------------------------------MESSAGE_READS----------------------------------------------------
func (r *MessageRepository) CreateMessageRead(ctx context.Context, messageRead *domain.MessageRead) (*domain.MessageRead, error) {
	if err := r.db.WithContext(ctx).Create(messageRead).Error; err != nil {
		return nil, err
	}
	return messageRead, nil
}

func (r *MessageRepository) GetMessageReadsByMessageID(ctx context.Context, id string) ([]domain.MessageRead, error) {
	var messageReads []domain.MessageRead
	if err := r.db.WithContext(ctx).Where("message_id = ?", id).Find(&messageReads).Error; err != nil {
		return nil, err
	}
	return messageReads, nil
}

func (r *MessageRepository) DeleteMessageRead(ctx context.Context, messageID, userID string) error {
	if err := r.db.WithContext(ctx).Where("message_id = $1 AND user_id = $2", messageID, userID).Delete(&domain.MessageRead{}).Error; err != nil {
		return err
	}
	return nil
}
