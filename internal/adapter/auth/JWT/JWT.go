package jwt

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/HellEaglee/Golang-Chat/internal/adapter/config"
	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
	"github.com/HellEaglee/Golang-Chat/internal/core/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTToken struct {
	secretKey        []byte
	duration         time.Duration
	refreshSecretKey []byte
	refreshDuration  time.Duration
	tr               port.TokenRepository
}

func New(config *config.Token, tokenRepo port.TokenRepository) (port.TokenService, error) {
	durationStr := config.Duration
	refreshDurationStr := config.DurationRefresh
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, util.ErrAccessTokenDuration
	}
	refreshDuration, err := time.ParseDuration(refreshDurationStr)
	if err != nil {
		return nil, util.ErrRefreshTokenDuration
	}
	secretKey := []byte(config.Secret)
	refreshSecretKey := []byte(config.SecretRefresh)

	return &JWTToken{
		secretKey:        secretKey,
		duration:         duration,
		refreshSecretKey: refreshSecretKey,
		refreshDuration:  refreshDuration,
		tr:               tokenRepo,
	}, nil
}

func (t *JWTToken) CreateToken(userID string, sessionID string) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", util.ErrInternal
	}

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(t.duration)
	claims := domain.TokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id.String(),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			Issuer:    "golang-chat",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.secretKey)
	if err != nil {
		return "", util.ErrAccessTokenCreation
	}
	return tokenString, nil
}

func (t *JWTToken) VerifyToken(tokenString string) (*domain.TokenPayload, error) {
	claims := &domain.TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return t.secretKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, util.ErrExpiredAccessToken
		}
		return nil, util.ErrInvalidAccessToken
	}
	if !token.Valid {
		return nil, util.ErrInvalidAccessToken
	}

	payload := &domain.TokenPayload{
		ID:        uuid.MustParse(claims.RegisteredClaims.ID),
		UserID:    claims.UserID,
		SessionID: claims.SessionID,
	}

	return payload, nil
}

func (t *JWTToken) CreateRefreshToken(ctx context.Context, userID string, sessionID string) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", util.ErrInternal
	}

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(t.refreshDuration)

	claims := domain.TokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id.String(),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			Issuer:    "golang-chat-refresh",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.refreshSecretKey)
	if err != nil {
		return "", util.ErrRefreshTokenCreation
	}

	tokenEntity := &domain.Token{
		ID:        id,
		UserID:    userID,
		SessionID: sessionID,
		Token:     tokenString,
		ExpiresAt: expiredAt,
	}

	_, err = t.tr.StoreRefreshToken(ctx, tokenEntity)
	if err != nil {
		return "", util.ErrInternal
	}

	return tokenString, nil
}

func (t *JWTToken) VerifyRefreshToken(ctx context.Context, tokenString string) (*domain.TokenPayload, error) {
	claims := &domain.TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return t.refreshSecretKey, nil
	})
	if err != nil {
		if err.Error() == "Token is expired" || err.Error() == "token is expired" || strings.Contains(err.Error(), "expired") {
			return nil, util.ErrExpiredAccessToken
		}
	}
	if !token.Valid {
		return nil, util.ErrInvalidAccessToken
	}
	storenToken, err := t.tr.GetTokenByID(ctx, claims.ID)
	if err != nil {
		return nil, util.ErrInvalidRefreshToken
	}

	if storenToken.SessionID != claims.SessionID {
		return nil, util.ErrInvalidRefreshToken
	}

	if time.Now().After(storenToken.ExpiresAt) {
		return nil, util.ErrExpiredRefreshToken
	}

	payload := &domain.TokenPayload{
		ID:        uuid.MustParse(claims.RegisteredClaims.ID),
		UserID:    claims.UserID,
		SessionID: claims.SessionID,
	}

	return payload, nil
}

func (t *JWTToken) ExtractClaimsFromToken(tokenString string) (*domain.TokenClaims, error) {
	claims := &domain.TokenClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return t.secretKey, nil
	})

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return nil, util.ErrInvalidAccessToken
	}

	return claims, nil
}

func (t *JWTToken) ExtractTokenID(tokenString string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", util.ErrInternal
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", util.ErrInternal
	}

	if jti, ok := claims["jti"].(string); ok {
		return jti, nil
	}

	return "", util.ErrInternal
}

func (t *JWTToken) GetTokenBySessionID(ctx context.Context, sessionID string) (*domain.Token, error) {
	return t.tr.GetTokenBySessionID(ctx, sessionID)
}

func (t *JWTToken) RevokeToken(ctx context.Context, tokenID string) error {
	return t.tr.RevokeToken(ctx, tokenID)
}

func (t *JWTToken) RefreshTokens(ctx context.Context, oldAccessToken, refreshToken string) (accessToken string, err error) {
	accessPayload, err := t.ExtractClaimsFromToken(oldAccessToken)
	if err != nil {
		return "", err
	}

	refreshPayload, err := t.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}
	if accessPayload.SessionID != refreshPayload.SessionID {
		return "", util.ErrInvalidSession
	}

	newAccessToken, err := t.CreateToken(refreshPayload.UserID, refreshPayload.SessionID)
	if err != nil {
		return "", err
	}

	err = t.RevokeToken(ctx, refreshPayload.ID.String())
	if err != nil {
		return "", err
	}
	_, err = t.CreateRefreshToken(ctx, refreshPayload.UserID, refreshPayload.SessionID)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}
