package service

import (
	. "brain-wave/internal/users-service/mo
	"context"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	. "brain-wave/internal/users-service/models"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) GetUserByID(ctx context.Context, userId int) (*User, error) {
	return s.repo.GetUserByID(ctx, userId)
}

func (s *Service) LoginUserByEmail(ctx context.Context, email, password string) (*User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	return user.User, nil
}

func (s *Service) LoginUserByUsername(ctx context.Context, username, password string) (*User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	return user.User, nil
}

func (s *Service) CreateUser(ctx context.Context, user *UserWithPassword) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	user.PasswordHash = string(hash)

	return s.repo.CreateUser(ctx, user)
}
