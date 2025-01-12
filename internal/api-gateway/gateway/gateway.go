package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"story-pulse/internal/api-gateway/resolver"
)

type ServiceOption struct {
	Name         string
	RegisterFunc func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
	DialOptions  []grpc.DialOption
}

type Wrapper struct {
	Path string
	Func func(http.Handler) http.Handler
}

type Gateway struct {
	logger *zap.SugaredLogger
	coons  []*grpc.ClientConn
}

func NewGateway(ctx context.Context, grpcMux *runtime.ServeMux, logger *zap.SugaredLogger, serviceOpts ...*ServiceOption) (*Gateway, error) {
	var gateway Gateway
	gateway.logger = logger

	for _, srvOpt := range serviceOpts {
		dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithResolvers(&resolver.Builder{})}
		dialOptions = append(dialOptions, srvOpt.DialOptions...)

		conn, err := dial(fmt.Sprintf("dynamic:///%s", srvOpt.Name), dialOptions...)
		if err != nil {
			return nil, err
		}

		gateway.coons = append(gateway.coons, conn)

		if err = srvOpt.RegisterFunc(ctx, grpcMux, conn); err != nil {
			return nil, err
		}
	}

	return &gateway, nil
}

func (g *Gateway) Close() error {
	var errs []error

	for _, conn := range g.coons {
		if err := conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func dial(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr, opts...)
}
