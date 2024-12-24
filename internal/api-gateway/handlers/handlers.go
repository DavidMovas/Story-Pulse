package handlers

import (
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	logger *zap.SugaredLogger
}

func NewHandler(logger *zap.SugaredLogger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	h.logger.Info("Alive")
	w.WriteHeader(http.StatusOK)
}
