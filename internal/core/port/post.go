package port

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
}

type PostService interface {
	CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
}
