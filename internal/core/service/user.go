package service

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
	"github.com/HellEaglee/Golang-Chat/internal/core/port"
	"github.com/HellEaglee/Golang-Chat/internal/core/util"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	name := util.EmailToName(user.Email)
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, util.ErrInternal
	}

	user.Name = name
	user.Password = hashedPassword
	return s.repo.CreateUser(ctx, user)
}
