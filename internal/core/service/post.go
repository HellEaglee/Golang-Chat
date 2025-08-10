package service

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
)

type PostService struct {
	repo port.PostRepository
}

func NewPostService(repo port.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	return s.repo.CreatePost(ctx, post)
}

func (s *PostService) GetPosts(ctx context.Context, skip string, limit string) ([]domain.Post, error) {
	return s.repo.GetPosts(ctx, skip, limit)
}
