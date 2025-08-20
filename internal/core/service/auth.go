package service

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
	"github.com/HellEaglee/Golang-Chat/internal/core/util"
	"github.com/google/uuid"
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

func (s *AuthService) Login(ctx context.Context, email, password string) (accessToken string, err error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == util.ErrDataNotFound {
			return "", util.ErrInvalidCredentials
		}
		return "", util.ErrInternal
	}

	sessionID := uuid.New().String()

	err = util.ComparePassword(password, user.Password)
	if err != nil {
		return "", util.ErrInvalidCredentials
	}

	accessToken, err = s.ts.CreateToken(user.ID.String(), sessionID)
	if err != nil {
		return "", err
	}

	_, err = s.ts.CreateRefreshToken(ctx, user.ID.String(), sessionID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *AuthService) Register(ctx context.Context, user *domain.User) (accessToken string, err error) {
	name := util.EmailToName(user.Email)
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return "", util.ErrInternal
	}

	user.Password = hashedPassword
	user.Name = name
	newUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", util.ErrInternal
	}

	sessionID := uuid.New().String()

	accessToken, err = s.ts.CreateToken(newUser.ID.String(), sessionID)
	if err != nil {
		return "", err
	}

	_, err = s.ts.CreateRefreshToken(ctx, newUser.ID.String(), sessionID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
