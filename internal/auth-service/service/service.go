package service

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"story-pulse/internal/auth-service/config"
	"story-pulse/internal/auth-service/repository"
	v1 "story-pulse/internal/shared/grpc/v1"
	"story-pulse/internal/shared/jwt"
)

type Service struct {
	users  v1.UsersServiceClient
	jwt    *jwt.Service
	repo   *repository.Repository
	logger *zap.SugaredLogger
}

func NewService(usersClient v1.UsersServiceClient, repo *repository.Repository, logger *zap.SugaredLogger, cfg *config.Config) *Service {
	jwtService := jwt.NewService(cfg.Secret, cfg.AccessExpirationTime, cfg.RefreshExpirationTime)

	return &Service{
		users:  usersClient,
		jwt:    jwtService,
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Register(ctx context.Context, req *v1.CreateUserRequest) (*v1.RegisterResponse, error) {
	res, err := s.users.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwt.GenerateToken(int(res.User.Id), res.User.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	refreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.repo.SaveRefreshToken(ctx, int(res.User.Id), res.User.Role, refreshToken, s.jwt.GetRefreshExpiration())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := &v1.RegisterResponse{
		User:         res.User,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	s.logger.Infow("Register success", "id", res.User.Id, "role", res.User.Role, "email", res.User.Email)
	s.logger.Infow("Access token", "token", accessToken)
	s.logger.Infow("Refresh token", "token", refreshToken)

	return result, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	tokeData, err := s.repo.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return s.jwt.GenerateToken(tokeData.UserID, tokeData.Role)
}
