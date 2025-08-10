package port

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
	// GetPostByID(ctx context.Context, id string) (*domain.Post, error)
	GetPosts(ctx context.Context, skip uint64, limit uint64) ([]domain.Post, error)
	// UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
	// DeletePost(ctx context.Context, id string) error
}

type PostService interface {
	CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
	// GetPost(ctx context.Context, id string) (*domain.Post, error)
	GetPosts(ctx context.Context, skip uint64, limit uint64) ([]domain.Post, error)
	// UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
	// DeletePost(ctx context.Context, id string) error
}
