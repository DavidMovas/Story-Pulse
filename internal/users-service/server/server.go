package server

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"net"
	"story-pulse/internal/shared/echox"
	"story-pulse/internal/shared/validation"
	"story-pulse/internal/users-service/config"
	"story-pulse/internal/users-service/handlers"
	"story-pulse/internal/users-service/repository"
	"story-pulse/internal/users-service/service"
	"time"
)

const (
	dbConnectionTime = 10 * time.Second
)

type Server struct {
	e      *echo.Echo
	logger *zap.SugaredLogger
	cfg    *config.Config
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	sugar := logger.Sugar()

	validation.SetupValidators()

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.HTTPErrorHandler = echox.ErrorHandler
	e.HideBanner = true
	e.HidePort = true

	db, err := connectDB(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	rep := repository.NewRepository(db)
	srv := service.NewService(rep)

	handler := handlers.NewHandler(srv, sugar)

	api := e.Group("/users")

	// Health check
	api.GET("/health", handler.Health)

	api.GET("/:userId", handler.GetUserByID)

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

func connectDB(ctx context.Context, connString string) (db *pgxpool.Pool, err error) {
	ticker := time.NewTicker(dbConnectionTime / 4)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("failed to connect to database: %w", ctx.Err())
		case <-ticker.C:
			db, err = pgxpool.New(ctx, connString)
			if err == nil {
				if err = db.Ping(ctx); err == nil {
					return db, nil
				}
			}
		}
	}
}
