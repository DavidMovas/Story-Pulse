package gateway

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryCookieGatewayInterceptor(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	tokens := md["refresh_token"]
	if len(tokens) == 0 {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	token := tokens[0]
	fmt.Printf("GATEWAY COOKIE INTER: %s\n", token)

	return invoker(ctx, method, req, reply, cc, opts...)
}
