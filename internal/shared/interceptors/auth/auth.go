package auth

import (
	v1 "brain-wave/internal/shared/grpc/v1"
	"go.uber.org/zap"
)

type Interceptor struct {
	client v1.AuthServiceClient
	logger *zap.SugaredLogger
	opts   map[string]*LevelOption
}

func NewAuthInterceptor(client v1.AuthServiceClient, pattern string, logger *zap.SugaredLogger, options []*LevelOption) *Interceptor {
	var opts = make(map[string]*LevelOption, len(options))
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
