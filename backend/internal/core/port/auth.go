package port

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
)

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, token *domain.Token) (*domain.Token, error)
	GetTokenByID(ctx context.Context, tokenID string) (*domain.Token, error)
	GetTokenBySessionID(ctx context.Context, sessionID string) (*domain.Token, error)
	RevokeToken(ctx context.Context, tokenID string) error
}

type TokenService interface {
	CreateToken(userID string, sessionID string) (string, error)
	CreateRefreshToken(ctx context.Context, userID string, sessionID string) (string, error)
	VerifyToken(token string) (*domain.TokenPayload, error)
	VerifyRefreshToken(ctx context.Context, tokenString string) (*domain.TokenPayload, error)
	ExtractClaimsFromToken(tokenString string) (*domain.TokenClaims, error)
	ExtractTokenID(tokenString string) (string, error)
	GetTokenBySessionID(ctx context.Context, sessionID string) (*domain.Token, error)
	RefreshTokens(ctx context.Context, oldAccessToken, oldRefreshToken string) (accessToken string, err error)
	RevokeToken(ctx context.Context, tokenID string) error
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (accessToken string, err error)
	Register(ctx context.Context, user *domain.User) (accessToken string, err error)
}
