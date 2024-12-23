package auth

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	v1 "story-pulse/internal/shared/grpc/v1"
)

type Interceptor struct {
	client v1.AuthServiceClient
	logger *zap.SugaredLogger
	opts   map[string]*AuthLevelOption
}

func NewAuthInterceptor(client v1.AuthServiceClient, pattern string, logger *zap.SugaredLogger, options []*AuthLevelOption) *Interceptor {
	var opts = make(map[string]*AuthLevelOption, len(options))
	for _, option := range options {
		opts[pattern+option.MethodName] = option
	}

	return &Interceptor{
		client: client,
		logger: logger,
		opts:   opts,
	}
}

// 1. Check if method has authLevel
// 2. If it is so try to take token from metadata
// 3. If there is no token return code
// 4. If there is token sending it to auth service

func UnaryServerInterceptor(a *Interceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		opt, ok := a.opts[info.FullMethod]
		if !ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "insufficient permission")
		}

		var token, userId string
		if tokens := md.Get("token"); len(tokens) > 0 {
			token = tokens[0]
			a.logger.Infof("TOKEN %s", token)
		}

		if userIds := md.Get("userid"); len(userIds) > 0 {
			userId = userIds[0]
			a.logger.Infof("USER_ID: %s", userId)
		}

		if a.client != nil {
			res, err := a.client.CheckToken(ctx, &v1.CheckTokenRequest{
				Token:  token,
				UserId: userId,
				Role:   opt.AuthLevel,
				Self:   opt.Self,
			})

			if err != nil {
				a.logger.Errorf("failed to request auth client: %v", err)
				return nil, status.Errorf(codes.Internal, "server error")
			}

			if res.Valid {
				return handler(ctx, req)
			}
		}

		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}
}
