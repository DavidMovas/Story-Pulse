package authentication

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	v1 "story-pulse/internal/shared/grpc/v1"
	"story-pulse/internal/shared/interceptors"
)

var _ interceptors.Interceptor = (*Interceptor)(nil)

type Interceptor struct {
	client v1.AuthServiceClient
	logger *zap.SugaredLogger
	opts   map[string]string
}

func NewAuthInterceptor(client v1.AuthServiceClient, pattern string, logger *zap.SugaredLogger, options []*AuthLevelOption) *Interceptor {
	var opts = make(map[string]string, len(options))
	for _, option := range options {
		opts[pattern+option.MethodName] = option.AuthLevel
		logger.Infow("Added AUTH option", "method", pattern+option.MethodName, "authLevel", option.AuthLevel)
	}

	return &Interceptor{
		client: client,
		opts:   opts,
	}
}

// 1. Check if method has authLevel
// 2. If it is so try to take token from metadata
// 3. If there is no token return code
// 4. If there is token sending it to auth service

func (a *Interceptor) Intercept(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	level, ok := a.opts[info.FullMethod]
	if !ok {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "insufficient permission")
	}

	tokenStr := md["authorization"][0]
	a.logger.Infof("TOKEN %s", tokenStr)

	res, err := a.client.CheckToken(ctx, &v1.CheckTokenRequest{
		Token: tokenStr,
		Role:  level,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "server error")
	}

	if res.Valid {
		return handler(ctx, req)
	}

	return nil, status.Errorf(codes.PermissionDenied, "permission denied")
}
