package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"story-pulse/internal/shared/consul"
	v1 "story-pulse/internal/shared/grpc/v1"
	net2 "story-pulse/internal/shared/net"
	"story-pulse/internal/users-service/config"
	"story-pulse/internal/users-service/handlers"
	"story-pulse/internal/users-service/repository"
	"story-pulse/internal/users-service/service"
)

type Server struct {
	grpcServer *grpc.Server
	httpServer *http.Server
	listener   net.Listener
	logger     *zap.SugaredLogger
	cfg        *config.Config

	closers []func() error
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	sugar := logger.Sugar().WithOptions(zap.WithCaller(false))

	db, err := connectDB(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	repo := repository.NewRepository(db)
	srv := service.NewService(repo)
	handler := handlers.NewHandler(srv, sugar)

	/*authClient, err := client.CreateServiceClient[v1.AuthServiceClient]("auth-service", v1.NewAuthServiceClient, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithResolvers(&resolver.Builder{}))
	if err != nil {
		sugar.Errorw("Failed to create auth service", "error", err)
		return nil, err
	}

	authInterceptor := auth.NewAuthInterceptor(authClient, "/userservice.v1.UsersService/", sugar, handler.GetAuthOptions())
	*/
	grpcServer := grpc.NewServer()

	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", handler.Health)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.WebPort),
		Handler: healthMux,
	}

	v1.RegisterUsersServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	return &Server{
		grpcServer: grpcServer,
		httpServer: httpServer,
		logger:     sugar,
		cfg:        cfg,
		closers: []func() error{func() error {
			db.Close()
			return nil
		}},
	}, nil
}

func (s *Server) Run() (err error) {
	s.listener, err = net2.ReservePort("", s.cfg.GRPCPort)
	if err != nil {
		return fmt.Errorf("cannot reserve grpc port: %w", err)
	}

	s.logger.Infof("starting server http port %d", s.cfg.WebPort)
	s.logger.Infof("starting server tcp port %d", s.cfg.GRPCPort)

	s.closers = append(s.closers, s.listener.Close)

	if err = s.register(); err != nil {
		return err
	}

	go func() {
		_ = s.httpServer.ListenAndServe()
	}()

	return s.grpcServer.Serve(s.listener)
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Infow("shutting down server")

	err := s.httpServer.Shutdown(ctx)
	s.grpcServer.GracefulStop()

	select {
	case <-ctx.Done():
		_ = s.httpServer.Close()
		s.grpcServer.Stop()
	}

	return withClosers(s.closers, err)
}

func (s *Server) Port() (int, error) {
	if s.listener == nil || s.listener.Addr() == nil {
		return 0, fmt.Errorf("server is not running")
	}

	return s.listener.Addr().(*net.TCPAddr).Port, nil
}

func (s *Server) register() error {
	consulCfg := api.DefaultConfig()
	consulCfg.Address = s.cfg.ConsulAddr
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", s.cfg.Address, s.cfg.WebPort),
		Interval:                       "10s",
		Timeout:                        "2s",
		DeregisterCriticalServiceAfter: "30s",
	}

	consulClient, err := api.NewClient(consulCfg)
	err = consul.RegisterService(consulClient, s.cfg.Name, s.cfg.Address, s.cfg.Tag, s.cfg.GRPCPort, check)
	if err != nil {
		return err
	}

	s.logger.Infow("Service registered in Consul", "Name", s.cfg.Name, "Address", s.cfg.Address, "Tag", s.cfg.Tag)
	return nil
}

func connectDB(ctx context.Context, connString string) (db *pgxpool.Pool, err error) {
	db, err = pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func withClosers(closers []func() error, err error) error {
	errs := []error{err}

	for i := len(closers) - 1; i >= 0; i-- {
		if err = closers[i](); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
