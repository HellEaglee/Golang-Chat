package port

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
)

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, token *domain.Token) (*domain.Token, error)
}

type TokenService interface {
	CreateToken(user *domain.User) (string, error)
	CreateRefreshToken(ctx context.Context, user *domain.User) (string, error)
	VerifyToken(token string) (*domain.TokenPayload, error)
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error)
	Register(ctx context.Context, user *domain.User) (accessToken string, refreshToken string, err error)
}
