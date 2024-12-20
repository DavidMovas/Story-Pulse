package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"story-pulse/internal/api-gateway/config"
	"story-pulse/internal/api-gateway/service"
	apperrors "story-pulse/internal/shared/error"
	"strings"
)

type Handler struct {
	service *service.Service
	cfg     *config.Config
	logger  *zap.SugaredLogger
}

func NewHandler(service *service.Service, cfg *config.Config, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		service: service,
		cfg:     cfg,
		logger:  logger,
	}
}

func (h *Handler) Health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func (h *Handler) Gateway(c echo.Context) error {
	parsedUrl, err := url.Parse(c.Request().URL.RequestURI())
	if err != nil {
		return c.JSON(http.StatusBadRequest, apperrors.BadRequest(err))
	}

	trimmed := strings.TrimPrefix(parsedUrl.String(), "/api/")
	servicePath := strings.SplitN(trimmed, "/", 2)[0]

	var serviceUrl string
	switch servicePath {
	case h.cfg.UsersService.ServicePath:
		serviceUrl = h.cfg.UsersService.ServiceURL
	case h.cfg.AuthService.ServicePath:
		serviceUrl = h.cfg.AuthService.ServiceURL
	case h.cfg.ContentService.ServicePath:
		serviceUrl = h.cfg.ContentService.ServiceURL
	case h.cfg.CommentService.ServicePath:
		serviceUrl = h.cfg.CommentService.ServiceURL
	case h.cfg.SearchService.ServicePath:
		serviceUrl = h.cfg.SearchService.ServiceURL
	default:
		return apperrors.BadRequest(fmt.Errorf("invalid path"))
	}

	res, err := h.service.Handle(serviceUrl, c.Request().Method, trimmed, parsedUrl.RawQuery, c.Request().Header, c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	for key, values := range res.Header {
		for _, value := range values {
			c.Response().Header().Add(key, value)
		}
	}

	c.Response().WriteHeader(res.StatusCode)
	_, err = io.Copy(c.Response().Writer, res.Body)
	return err
}
