package handlers

import (
	"context"
	"go.uber.org/zap"
	grpc "story-pulse/internal/shared/grpc/v1"
	. "story-pulse/internal/users-service/models"
	. "story-pulse/internal/users-service/service"
)

var _ grpc.UsersServiceServer = (*Handler)(nil)

type Handler struct {
	service *Service
	logger  *zap.SugaredLogger

	grpc.UnimplementedUsersServiceServer
}

func NewHandler(service *Service, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) GetUserByID(ctx context.Context, request *grpc.GetUserByIDRequest) (*grpc.GetUserByIDResponse, error) {
	userId := int(request.GetId())
	user, err := h.service.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &grpc.GetUserByIDResponse{User: user.ToGRPC()}, nil
}

func (h *Handler) CreateUser(ctx context.Context, request *grpc.CreateUserRequest) (*grpc.CreateUserResponse, error) {
	user, err := h.service.CreateUser(ctx, ToUserWithPassword(request))
	if err != nil {
		return nil, err
	}

	return &grpc.CreateUserResponse{User: user.ToGRPC()}, nil
}
