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

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) GetUsers(ctx context.Context, skip uint64, limit uint64) ([]domain.User, error) {
	return s.repo.GetUsers(ctx, skip, limit)
}

func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, util.ErrInternal
	}
	user.Password = hashedPassword
	return s.repo.UpdateUser(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.DeleteUser(ctx, id)
}
