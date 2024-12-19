package handlers

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"story-pulse/internal/shared/echox"
	apperrors "story-pulse/internal/shared/error"
	. "story-pulse/internal/users-service/models"
	. "story-pulse/internal/users-service/service"
)

type Handler struct {
	service *Service
	logger  *zap.SugaredLogger
}

func NewHandler(service *Service, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func (h *Handler) GetUserByID(c echo.Context) error {
	ctx := c.Request().Context()
	req, err := echox.BindAndValidate[GetUserByIDRequest](c)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	user, err := h.service.GetUserByID(ctx, req.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()
	req, err := echox.BindAndValidate[CreateUserRequest](c)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	user, err := h.service.CreateUser(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}
