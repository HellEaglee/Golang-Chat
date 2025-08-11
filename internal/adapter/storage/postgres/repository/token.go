package repository

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/adapter/storage/postgres"
	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
)

type TokenRepository struct {
	db *postgres.DB
}

func NewTokenRepository(db *postgres.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) StoreRefreshToken(ctx context.Context, token *domain.Token) (*domain.Token, error) {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return nil, err
	}
	return token, nil
}
