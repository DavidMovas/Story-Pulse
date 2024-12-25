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
	"story-pulse/internal/shared/resolver"
)

type ServiceOption struct {
	Name         string
	MuxOptions   []runtime.ServeMuxOption
	RegisterFunc func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
	DialOptions  []grpc.DialOption
}

type Gateway struct {
	mux    *http.ServeMux
	logger *zap.SugaredLogger
	coons  []*grpc.ClientConn
}

func NewGateway(ctx context.Context, httpMux *http.ServeMux, logger *zap.SugaredLogger, globalMuxOptions []runtime.ServeMuxOption, serviceOpts ...*ServiceOption) (*Gateway, error) {
	var gateway Gateway
	gateway.logger = logger
	gateway.mux = httpMux

	for _, srvOpt := range serviceOpts {
		dialOptions := []grpc.DialOption{grpc.WithResolvers(&resolver.Builder{}), grpc.WithTransportCredentials(insecure.NewCredentials())}
		dialOptions = append(dialOptions, srvOpt.DialOptions...)

		conn, err := dial(fmt.Sprintf("dynamic:///%s", srvOpt.Name), dialOptions...)
		if err != nil {
			return nil, err
		}

		gateway.coons = append(gateway.coons, conn)

		var options []runtime.ServeMuxOption
		options = append(options, globalMuxOptions...)
		options = append(options, srvOpt.MuxOptions...)

		mux := runtime.NewServeMux(options...)
		if err = srvOpt.RegisterFunc(ctx, mux, conn); err != nil {
			return nil, err
		}

		httpMux.Handle("/", mux)
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
