package handlers

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
	"story-pulse/internal/auth-service/service"
	v1 "story-pulse/internal/shared/grpc/v1"
	"story-pulse/internal/shared/validation"
)

var _ v1.AuthServiceServer = (*Handler)(nil)

type Handler struct {
	service *service.Service
	logger  *zap.SugaredLogger

	v1.UnimplementedAuthServiceServer
}

func NewHandler(service *service.Service, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) RegisterUser(ctx context.Context, request *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	if request.Email == "" && request.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "email and username required")
	}

	if err := validation.Validate("email", request.Email, true); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := validation.Validate("username", request.Username, true); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := validation.Validate("password", request.Password, true); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	res, err := h.service.Register(ctx, &v1.CreateUserRequest{
		Email:    request.Email,
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}

	return &v1.RegisterResponse{
		User:         res.User,
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (h *Handler) LoginUser(ctx context.Context, request *v1.LoginRequest) (*v1.LoginResponse, error) {
	if request.Email == nil && request.Username == nil {
		return nil, status.Error(codes.InvalidArgument, "email or username required")
	}

	if err := validation.Validate("email", request.Email, false); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := validation.Validate("username", request.Username, false); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := validation.Validate("password", request.Password, true); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &v1.LoginResponse{}, nil
}

func (h *Handler) RefreshToken(ctx context.Context, request *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	data := md.Get("refresh_token")
	if len(data) == 0 {
		return nil, status.Error(codes.Unauthenticated, "refresh token is empty")
	}

	refreshToken := data[0]
	h.logger.Info("RefreshToken called", zap.String("refresh_token", refreshToken))

	accessToken, err := h.service.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &v1.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (h *Handler) Health(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("ok"))
}
