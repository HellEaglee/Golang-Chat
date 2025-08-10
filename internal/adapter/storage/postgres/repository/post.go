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

func (r *PostRepository) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	var post domain.Post
	if err := r.db.WithContext(ctx).First(post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) GetPosts(ctx context.Context, skip string, limit string) ([]domain.Post, error) {
	var posts []domain.Post
	if err := r.db.WithContext(ctx).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
