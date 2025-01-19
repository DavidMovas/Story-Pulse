package handlers

import (
	grpc "brain-wave/internal/shared/grpc/v1"
	auth "brain-wave/internal/shared/interceptors/auth"
	. "brain-wave/internal/users-service/models"
	. "brain-wave/internal/users-service/service"
	"context"
	"go.uber.org/zap"
	"net/http"
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
			{"GetUserByID", "user", false},
			{"CreateUser", "admin", false},
		},
	}
}

func (h *Handler) GetUserByID(ctx context.Context, request *grpc.GetUserByIDRequest) (*grpc.GetUserByIDResponse, error) {
	userId := int(request.Id)
	user, err := h.service.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	h.logger.Info("GetUserByID invoked")

	return &grpc.GetUserByIDResponse{User: user.ToGRPC()}, nil
}

func (h *Handler) LoginUserByEmail(ctx context.Context, request *grpc.LoginUserByEmailRequest) (*grpc.LoginUserResponse, error) {
	user, err := h.service.LoginUserByEmail(ctx, request.Email, request.Password)
	if err != nil {
		return nil, err
	}

	return &grpc.LoginUserResponse{User: user.ToGRPC()}, nil
}

func (h *Handler) LoginUserByUsername(ctx context.Context, request *grpc.LoginUserByUsernameRequest) (*grpc.LoginUserResponse, error) {
	user, err := h.service.LoginUserByUsername(ctx, request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	return &grpc.LoginUserResponse{User: user.ToGRPC()}, nil
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
