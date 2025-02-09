package server

import (
	"brain-wave/internal/api-gateway/config"
	"brain-wave/internal/api-gateway/gateway"
	"brain-wave/internal/api-gateway/handlers"
	"brain-wave/internal/api-gateway/middlewares"
	"brain-wave/internal/api-gateway/mux"
	"brain-wave/internal/api-gateway/options"
	v1 "brain-wave/internal/shared/grpc/v1"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
)

type Server struct {
	gateway *gateway.Gateway
	mux     *chi.Mux
	logger  *zap.SugaredLogger
	cfg     *config.Config

	ctx    context.Context
	cancel context.CancelFunc
}

func NewServer(cfg *config.Config) (*Server, error) {
	serverCtx, cancel := context.WithCancel(context.Background())
	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	sugar := logger.Sugar().WithOptions(zap.WithCaller(false))
	grpcMuxOpts := []runtime.ServeMuxOption{runtime.WithErrorHandler(options.CustomErrorHandler)}

	// GRPC Mux
	grpcMux := runtime.NewServeMux(grpcMuxOpts...)

	// HTTP Mux
	httpMux := chi.NewMux()
	httpMux.Use(middlewares.NewLoggerMiddleware(sugar))

	// Register middlewares and routes
	mux.Register(httpMux, grpcMux)

	handler := handlers.NewHandler(sugar)
	httpMux.HandleFunc("/health", handler.Health)

	serviceOpts := []*gateway.ServiceOption{
		{
			Name:         cfg.UsersService.ServicePath,
			RegisterFunc: v1.RegisterUsersServiceHandler,
			DialOptions: []grpc.DialOption{
				grpc.WithPerRPCCredentials(options.NewAuthenticateCredentials()),
			},
		},
		{
			Name:         cfg.AuthService.ServicePath,
			RegisterFunc: v1.RegisterAuthServiceHandler,
		},
	}

	gt, err := gateway.NewGateway(serverCtx, grpcMux, sugar, serviceOpts...)
	if err != nil {
		cancel()
		return nil, err
	}

	return &Server{
		gateway: gt,
		mux:     httpMux,
		ctx:     serverCtx,
		logger:  sugar,
		cfg:     cfg,
		cancel:  cancel,
	}, nil
}

func (s *Server) Run() error {
	port := s.cfg.WebPort
	s.logger.Infof("starting server on port %d", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), s.mux)
}

func (s *Server) Stop() error {
	s.cancel()
	return s.gateway.Close()
}
