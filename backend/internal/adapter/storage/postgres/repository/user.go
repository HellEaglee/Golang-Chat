package repository

import (
	"context"

	"github.com/HellEaglee/Golang-Chat/internal/adapter/storage/postgres"
	"github.com/HellEaglee/Golang-Chat/internal/core/domain"
)

type UserRepository struct {
	db *postgres.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUsers(ctx context.Context, skip uint64, limit uint64) ([]domain.User, error) {
	var users []domain.User
	if err := r.db.WithContext(ctx).Limit(int(limit)).Offset(int(skip)).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	var updatedUser domain.User
	query := `UPDATE users SET name = $2, email = $3, password = $4, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL RETURNING *`

	if err := r.db.WithContext(ctx).Raw(query, user.ID, user.Name, user.Email, user.Password).Scan(&updatedUser).Error; err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.User{}).Error; err != nil {
		return err
	}
	return nil
}
