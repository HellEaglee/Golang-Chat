package repository

import (
	"context"
	"time"

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

func (r *TokenRepository) GetTokenByID(ctx context.Context, tokenID string) (*domain.Token, error) {
	var token domain.Token
	if err := r.db.WithContext(ctx).Where("id = ? AND revoked_at IS NULL", tokenID).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokenRepository) RevokeToken(ctx context.Context, tokenID string) error {
	return r.db.WithContext(ctx).Model(&domain.Token{}).Where("id = ?", tokenID).Update("revoked_at", time.Now()).Error
}
