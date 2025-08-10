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

func (s *PostService) GetPost(ctx context.Context, id string) (*domain.Post, error) {
	return s.repo.GetPostByID(ctx, id)
}

func (s *PostService) GetPosts(ctx context.Context, skip uint64, limit uint64) ([]domain.Post, error) {
	return s.repo.GetPosts(ctx, skip, limit)
}

func (s *PostService) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	return s.repo.UpdatePost(ctx, post)
}

func (s *PostService) DeletePost(ctx context.Context, id string) error {
	return s.repo.DeletePost(ctx, id)
}
