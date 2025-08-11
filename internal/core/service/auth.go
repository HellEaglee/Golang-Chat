package service

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
	"github.com/HellEaglee/Golang-Chat/internal/core/util"
)

type AuthService struct {
	repo port.UserRepository
	ts   port.TokenService
}

func NewAuthService(repo port.UserRepository, ts port.TokenService) *AuthService {
	return &AuthService{
		repo: repo, ts: ts,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == util.ErrDataNotFound {
			return "", "", util.ErrInvalidCredentials
		}
		return "", "", util.ErrInternal
	}

	err = util.ComparePassword(password, user.Password)
	if err != nil {
		return "", "", util.ErrInvalidCredentials
	}

	accessToken, err = s.ts.CreateToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.ts.CreateRefreshToken(ctx, user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Register(ctx context.Context, user *domain.User) (accessToken string, refreshToken string, err error) {
	name := util.EmailToName(user.Email)
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return "", "", util.ErrInternal
	}

	user.Password = hashedPassword
	user.Name = name
	newUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", "", util.ErrInternal
	}

	accessToken, err = s.ts.CreateToken(newUser)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.ts.CreateRefreshToken(ctx, newUser)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, oldAccessToken, oldRefreshToken string) (accessToken, refreshToken string, err error) {
	_, err = s.ts.VerifyToken(oldAccessToken)
	if err != nil {
		return "", "", err
	}

	payload, err := s.ts.VerifyRefreshToken(ctx, oldRefreshToken)
	if err != nil {
		return "", "", err
	}

	tokenID, err := s.ts.ExtractTokenID(oldRefreshToken)
	if err != nil {
		return "", "", err
	}

	user := &domain.User{
		ID: payload.UserID,
	}

	newAccessToken, err := s.ts.CreateToken(user)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.ts.CreateRefreshToken(ctx, user)
	if err != nil {
		return "", "", err
	}

	err = s.ts.RevokeToken(ctx, tokenID)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
