package service

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
)

type ChatService struct {
	repo port.ChatRepository
}

func NewChatService(repo port.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

// ----------------------------------------------------CHATS----------------------------------------------------
func (s *ChatService) CreateChat(ctx context.Context, chat *domain.Chat) (*domain.Chat, error) {
	return s.repo.CreateChat(ctx, chat)
}

func (s *ChatService) GetChatByID(ctx context.Context, id string) (*domain.Chat, error) {
	return s.repo.GetChatByID(ctx, id)
}

func (s *ChatService) GetChatsByUserID(ctx context.Context, id string) ([]domain.Chat, error) {
	return s.repo.GetChatsByUserID(ctx, id)
}

func (s *ChatService) GetChats(ctx context.Context, skip uint64, limit uint64) ([]domain.Chat, error) {
	return s.repo.GetChats(ctx, skip, limit)
}

func (s *ChatService) UpdateChat(ctx context.Context, chat *domain.Chat) (*domain.Chat, error) {
	return s.repo.UpdateChat(ctx, chat)
}

func (s *ChatService) DeleteChat(ctx context.Context, id string) error {
	return s.repo.DeleteChat(ctx, id)
}

// ----------------------------------------------------CHAT_PARTICIPANTS----------------------------------------------------
func (s *ChatService) CreateChatParticipant(ctx context.Context, chatParticipant *domain.ChatParticipant) (*domain.ChatParticipant, error) {
	return s.repo.CreateChatParticipant(ctx, chatParticipant)
}

func (s *ChatService) GetChatParticipantByChatIDUserID(ctx context.Context, chatID, userID string) (*domain.ChatParticipant, error) {
	return s.repo.GetChatParticipantByChatIDUserID(ctx, chatID, userID)
}

func (s *ChatService) GetChatParticipantsByChatID(ctx context.Context, id string) ([]domain.ChatParticipant, error) {
	return s.repo.GetChatParticipantsByChatID(ctx, id)
}

func (s *ChatService) UpdateChatParticipant(ctx context.Context, chatParticipant *domain.ChatParticipant) (*domain.ChatParticipant, error) {
	return s.repo.UpdateChatParticipant(ctx, chatParticipant)
}

func (s *ChatService) DeleteChatParticipant(ctx context.Context, chatID, userID string) error {
	return s.repo.DeleteChatParticipant(ctx, chatID, userID)
}
