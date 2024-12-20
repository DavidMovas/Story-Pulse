package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
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
	requestUrl := c.Request().URL.RequestURI()
	parsedUrl, err := url.Parse(requestUrl)
	if err != nil {
		return c.JSON(http.StatusBadRequest, apperrors.BadRequest(err))
	}

	method := c.Request().Method
	path := strings.TrimPrefix(parsedUrl.Path, "/")
	query := parsedUrl.RawQuery
	body := c.Request().Body
	defer func() {
		_ = body.Close()
	}()

	fmt.Printf("FULL REQUEST URL: %s\n\n", requestUrl)
	fmt.Printf("PATH: %s\n", path)
	fmt.Printf("QUERY: %s\n", query)

	var serviceUrl string
	switch path {
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
		return c.JSON(http.StatusBadRequest, apperrors.BadRequest(fmt.Errorf("invalid path")))
	}

	res, err := h.service.Handle(serviceUrl, method, path, query, body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(res.StatusCode, res.Body)
}
