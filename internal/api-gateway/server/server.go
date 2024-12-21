package server

import (
	"context"
	"fmt"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"net/http"
	"story-pulse/internal/api-gateway/config"
	"story-pulse/internal/api-gateway/gateway"
	"story-pulse/internal/api-gateway/handlers"
	"story-pulse/internal/api-gateway/middlewares"
	"story-pulse/internal/api-gateway/options"
	grpcServices "story-pulse/internal/shared/grpc/v1"
)

type Server struct {
	gateway *gateway.Gateway
	mux     *http.ServeMux
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

	mux := http.NewServeMux()

	handler := handlers.NewHandler(sugar)
	mux.HandleFunc("/health", handler.Health)

	muxOpts := []gwruntime.ServeMuxOption{
		gwruntime.WithErrorHandler(options.CustomErrorHandler),
		gwruntime.WithMiddlewares(middlewares.NewLoggerMiddleware(sugar)),
	}

	serviceOpts := []*gateway.ServiceOption{
		{
			Url:          cfg.UsersService.ServiceURL,
			RegisterFunc: grpcServices.RegisterUsersServiceHandler,
		},
	}

	gt, err := gateway.NewGateway(serverCtx, sugar, muxOpts, serviceOpts...)
	if err != nil {
		cancel()
		return nil, err
	}

	mux.Handle("/", gt.Proxy())

	return &Server{
		gateway: gt,
		mux:     mux,
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
