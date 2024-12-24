package gateway

import (
	"context"
	"errors"
	"fmt"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

type ServiceOption struct {
	Name         string
	RegisterFunc func(ctx context.Context, mux *gwruntime.ServeMux, conn *grpc.ClientConn) error
	Options      []grpc.DialOption
}

type Gateway struct {
	mux    *gwruntime.ServeMux
	logger *zap.SugaredLogger
	coons  []*grpc.ClientConn
}

func NewGateway(ctx context.Context, logger *zap.SugaredLogger, opts []gwruntime.ServeMuxOption, serviceOpts ...*ServiceOption) (*Gateway, error) {
	mux := gwruntime.NewServeMux(opts...)

	var gateway Gateway
	gateway.logger = logger
	gateway.mux = mux

	for _, opt := range serviceOpts {
		dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		dialOptions = append(dialOptions, opt.Options...)

		conn, err := dial(fmt.Sprintf("dynamic:///%s", opt.Name), dialOptions...)
		if err != nil {
			return nil, err
		}

		gateway.coons = append(gateway.coons, conn)

		if err = opt.RegisterFunc(ctx, mux, conn); err != nil {
			return nil, err
		}
	}

	return &gateway, nil
}

func (g *Gateway) Proxy() http.Handler {
	return g.mux
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
