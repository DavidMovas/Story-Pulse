package service

import (
	"brain-wave/internal/auth-service/config"
	"brain-wave/internal/auth-service/repository"
	v1 "brain-wave/internal/shared/grpc/v1"
	"brain-wave/internal/shared/jwt"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc
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

	return result, nil
}

func (s *Service) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	var res *v1.LoginUserResponse
	var err error
	if req.Email != nil {
		res, err = s.users.LoginUserByEmail(ctx, &v1.LoginUserByEmailRequest{Email: *req.Email, Password: req.Password})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else if req.Username != nil {
		res, err = s.users.LoginUserByUsername(ctx, &v1.LoginUserByUsernameRequest{Username: *req.Username, Password: req.Password})
	}

	if res == nil {
		return nil, status.Error(codes.NotFound, "User not found")
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

	result := &v1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         res.User,
	}

	return result, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	tokeData, err := s.repo.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return s.jwt.GenerateToken(tokeData.UserID, tokeData.Role)
}
