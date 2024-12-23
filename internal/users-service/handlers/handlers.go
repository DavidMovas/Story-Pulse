package handlers

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	grpc "story-pulse/internal/shared/grpc/v1"
	auth "story-pulse/internal/shared/interceptors/authentication"
	. "story-pulse/internal/users-service/models"
	. "story-pulse/internal/users-service/service"
)

var _ grpc.UsersServiceServer = (*Handler)(nil)

type Handler struct {
	service       *Service
	logger        *zap.SugaredLogger
	authLevelOpts []*auth.AuthLevelOption

	grpc.UnimplementedUsersServiceServer
}

func NewHandler(service *Service, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
		authLevelOpts: []*auth.AuthLevelOption{
			{"GetUserByID", "user"},
			{"CreateUser", "admin"},
		},
	}
}

func (h *Handler) GetUserByID(ctx context.Context, request *grpc.GetUserByIDRequest) (*grpc.GetUserByIDResponse, error) {
	userId := int(request.GetId())
	user, err := h.service.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	h.logger.Info("GetUserByID invoked")

	return &grpc.GetUserByIDResponse{User: user.ToGRPC()}, nil
}

func (h *Handler) CreateUser(ctx context.Context, request *grpc.CreateUserRequest) (*grpc.CreateUserResponse, error) {
	user, err := h.service.CreateUser(ctx, ToUserWithPassword(request))
	if err != nil {
		return nil, err
	}

	return &grpc.CreateUserResponse{User: user.ToGRPC()}, nil
}

func (h *Handler) Health(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("ok"))
}

func (h *Handler) GetAuthOptions() []*auth.AuthLevelOption {
	return h.authLevelOpts
}
