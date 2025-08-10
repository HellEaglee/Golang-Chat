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
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) GetPosts(ctx context.Context, skip uint64, limit uint64) ([]domain.Post, error) {
	var posts []domain.Post
	if err := r.db.WithContext(ctx).Limit(int(limit)).Offset(int(skip)).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	var updatedPost domain.Post

	query := `UPDATE posts SET title = $2, description = $3, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL RETURNING *`

	if err := r.db.WithContext(ctx).Raw(query, post.ID, post.Title, post.Description).Scan(&updatedPost).Error; err != nil {
		return nil, err
	}
	return &updatedPost, nil
}

// func (r *PostRepository) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
// 	if err := r.db.WithContext(ctx).Model(&domain.Post{}).Where("id = ? AND deleted_at IS NULL", post.ID).Updates(domain.Post{Title: post.Title, Description: post.Description}).Error; err != nil {
// 		return nil, err
// 	}
// 	var updatedPost domain.Post
// 	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", post.ID).First(&updatedPost).Error; err != nil {
// 		return nil, err
// 	}
// 	return &updatedPost, nil
// }
