package resolver

import (
	"github.com/hashicorp/consul/api"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
)

func init() { resolver.Register(&resolverBuilder{}) }

const (
	tickerTimeout = time.Second * 5
	consulAddress = "127.0.0.1:8500"
)

var _ resolver.Builder = (*resolverBuilder)(nil)
var _ resolver.Resolver = (*Resolver)(nil)

type resolverBuilder struct{}

func (*resolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r, err := NewResolver(target, cc, opts)
	if err != nil {
		return nil, err
	}

	go r.watchUpdates()
	return r, nil
}

func (*resolverBuilder) Scheme() string {
	return "dynamic"
}

type Resolver struct {
	consul    *api.Client
	target    resolver.Target
	cc        resolver.ClientConn
	opts      resolver.BuildOptions
	addresses []resolver.Address
	mu        sync.RWMutex

	ticker        *time.Ticker
	tickerTimeout time.Duration
}

func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {
	//Async resolving
}

func (r *Resolver) Close() {
	r.ticker.Stop()
}

func NewResolver(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (*Resolver, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulAddress

	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	var r Resolver
	r.consul = client
	r.target = target
	r.cc = cc
	r.opts = opts
	r.addresses = make([]resolver.Address, 0)
	r.tickerTimeout = tickerTimeout
	r.ticker = time.NewTicker(tickerTimeout)

	return &r, nil
}

func (r *Resolver) watchUpdates() {
	for {
		select {
		case <-r.ticker.C:
			services, _, err := r.consul.Health().Service(r.target.Endpoint(), "", true, nil)
			if err != nil {
				log.Errorf("failed to load healty services from consul: %v", err)
				r.ticker.Reset(r.tickerTimeout)
				continue
			}

			var addresses []resolver.Address
			for _, service := range services {
				addr := service.Service.Address
				if addr == "" {
					addr = service.Node.Address
				}

				addresses = append(addresses, resolver.Address{Addr: addr + ":" + string(rune(service.Service.Port))})
			}

			r.mu.Lock()
			r.addresses = addresses
			r.mu.Unlock()

			if err = r.cc.UpdateState(resolver.State{Addresses: addresses}); err != nil {
				log.Errorf("failed to update resolver state: %v", err)
			}
		}
	}
}
