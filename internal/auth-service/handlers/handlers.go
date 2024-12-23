package handlers

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	v1 "story-pulse/internal/shared/grpc/v1"
)

var _ v1.AuthServiceServer = (*Handler)(nil)

type Handler struct {
	logger *zap.SugaredLogger

	v1.UnimplementedAuthServiceServer
}

func NewHandler(logger *zap.SugaredLogger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) Health(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("ok"))
}

func (h *Handler) GenerateToken(_ context.Context, _ *v1.GenerateTokenRequest) (*v1.GenerateTokenResponse, error) {
	return &v1.GenerateTokenResponse{}, nil
}

func (h *Handler) CheckToken(_ context.Context, request *v1.CheckTokenRequest) (*v1.CheckTokenResponse, error) {
	h.logger.Infow("Check Token", "token", request.Token, "role", request.Role, "userId", request.UserId, "self", request.Self)

	return &v1.CheckTokenResponse{
		Valid: false,
	}, nil
}
