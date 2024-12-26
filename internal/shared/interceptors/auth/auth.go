package auth

import (
	"go.uber.org/zap"
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
