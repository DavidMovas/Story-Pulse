package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	apperrors "story-pulse/internal/shared/error"
	. "story-pulse/internal/users-service/models"
	. "story-pulse/internal/users-service/repository"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) GetUserByID(ctx context.Context, userId int) (*User, error) {
	return s.repo.GetUserByID(ctx, userId)
}

func (s *Service) CreateUser(ctx context.Context, user *CreateUserRequest) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	return s.repo.CreateUser(ctx, user.ToUserWithPassword(string(hash)))
}