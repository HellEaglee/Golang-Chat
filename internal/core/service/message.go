package service

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
)

type MessageService struct {
	repo port.MessageRepository
}

func NewMessageService(repo port.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

// ----------------------------------------------------MESSAGES----------------------------------------------------
func (s *MessageService) CreateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error) {
	return s.repo.CreateMessage(ctx, message)
}

func (s *MessageService) GetMessage(ctx context.Context, id string) (*domain.Message, error) {
	return s.repo.GetMessageByID(ctx, id)
}

func (s *MessageService) GetMessagesByChatID(ctx context.Context, chatID string) ([]domain.Message, error) {
	return s.repo.GetMessagesByChatID(ctx, chatID)
}

func (s *MessageService) UpdateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error) {
	return s.repo.UpdateMessage(ctx, message)
}

func (s *MessageService) DeleteMessage(ctx context.Context, id string) error {
	return s.repo.DeleteMessage(ctx, id)
}

// ----------------------------------------------------MESSAGE_READS----------------------------------------------------
func (s *MessageService) CreateMessageRead(ctx context.Context, messageRead *domain.MessageRead) (*domain.MessageRead, error) {
	return s.repo.CreateMessageRead(ctx, messageRead)
}

func (s *MessageService) GetMessageReadsByMessageID(ctx context.Context, id string) ([]domain.MessageRead, error) {
	return s.repo.GetMessageReadsByMessageID(ctx, id)
}

func (s *MessageService) DeleteMessageRead(ctx context.Context, messageID, userID string) error {
	return s.repo.DeleteMessageRead(ctx, messageID, userID)
}
