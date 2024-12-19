package handlers

import (
	"github.com/labstack/echo/v4"
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

func (h *Handler) Test(c echo.Context) error {
	h.logger.Infof("Calling test handler")
	return c.JSON(http.StatusOK, "Hello World")
}

func (h *Handler) Health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
