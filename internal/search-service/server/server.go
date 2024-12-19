package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"net"
	"story-pulse/internal/search-service/config"
	"story-pulse/internal/search-service/handlers"
)

type Server struct {
	e      *echo.Echo
	logger *zap.SugaredLogger
	cfg    *config.Config
}

func NewServer(cfg *config.Config) (*Server, error) {
	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	sugar := logger.Sugar()

	e := echo.New()

	e.Use(middleware.Recover())
	e.HideBanner = true
	e.HidePort = true

	handler := handlers.NewHandler(sugar)

	api := e.Group("/search-service")
	api.GET("/health", handler.Health)

	return &Server{
		e:      e,
		logger: sugar,
		cfg:    cfg,
	}, nil
}

func (s *Server) Run() error {
	port := s.cfg.WebPort
	s.logger.Infof("starting server on port %d", port)

	return s.e.Start(fmt.Sprintf(":%d", port))
}

func (s *Server) Stop(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) Port() (int, error) {
	listener := s.e.Listener
	if listener == nil {
		return 0, fmt.Errorf("no listener configured")
	}

	addr := listener.Addr()
	if addr == nil {
		return 0, fmt.Errorf("no listener address")
	}

	return addr.(*net.TCPAddr).Port, nil
}
