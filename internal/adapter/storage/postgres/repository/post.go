package repository

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/adapter/storage/postgres"
	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
)

type PostRepository struct {
	db *postgres.DB
}

func NewPostRepository(db *postgres.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	if err := r.db.WithContext(ctx).Create(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}
