package logging

import (
	"context"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"time"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Infof("INTERCEPT METHOD %s INVOKED AT: %s\n", info.FullMethod, time.Now().String())
		return handler(ctx, req)
	}
}
