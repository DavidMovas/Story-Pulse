package service

import (
	"context"
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
