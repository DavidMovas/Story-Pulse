package server

import (
	"context"
	"errors"
	"fmt"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/http"
	"story-pulse/internal/api-gateway/config"
	"story-pulse/internal/api-gateway/middlewares"
	grpcServices "story-pulse/internal/shared/grpc/v1"
)

type Server struct {
	mx     *gwruntime.ServeMux
	logger *zap.SugaredLogger
	cfg    *config.Config

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

	mx := gwruntime.NewServeMux(
		gwruntime.WithMiddlewares(middlewares.NewLoggerMiddleware(sugar)),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))}

	var errorsArray []error
	for _, serviceUrl := range []string{cfg.UsersService.ServiceURL, cfg.AuthService.ServiceURL} {
		err := grpcServices.RegisterUsersServiceHandlerFromEndpoint(serverCtx, mx, serviceUrl, opts)
		sugar.Infow("connecting to grpc service", "service_url", serviceUrl)
		if err != nil {
			errorsArray = append(errorsArray, err)
		}
	}

	if len(errorsArray) > 0 {
		cancel()
		return nil, errors.Join()
	}

	return &Server{
		mx:     mx,
		ctx:    serverCtx,
		logger: sugar,
		cfg:    cfg,
		cancel: cancel,
	}, nil
}

func (s *Server) Run() error {
	port := s.cfg.WebPort
	s.logger.Infof("starting server on port %d", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), s.mx)
}

func (s *Server) Stop() { s.cancel() }
