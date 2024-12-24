package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *Service) CreateUser(ctx context.Context, user *UserWithPassword) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	user.PasswordHash = string(hash)

	return s.repo.CreateUser(ctx, user)
}
