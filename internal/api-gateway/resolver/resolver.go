package resolver

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
)

func init() { resolver.Register(&resolverBuilder{}) }

const (
	tickerTimeout = time.Second * 15
	consulAddress = "http://consul:8500"
)

var _ resolver.Builder = (*resolverBuilder)(nil)
var _ resolver.Resolver = (*Resolver)(nil)

type resolverBuilder struct{}

func (*resolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r, err := NewResolver(target, cc, opts)
	if err != nil {
		return nil, err
	}

	r.refreshAddresses()
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
	mu        sync.Mutex

	ticker        *time.Ticker
	tickerTimeout time.Duration
}

func (r *Resolver) ResolveNow(_ resolver.ResolveNowOptions) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.addresses) == 0 {
		r.refreshAddresses()
	}

	_ = r.UpdateState(resolver.State{Addresses: r.addresses})
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
				log.Errorf("failed to load healty services from consul with endpoint: %s error: %v", r.target.Endpoint(), err)
				r.ticker.Reset(r.tickerTimeout)
				continue
			}

			var addresses = make([]resolver.Address, len(services))
			for i, service := range services {
				addresses[i] = resolver.Address{
					Addr: fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port),
				}
			}

			r.mu.Lock()
			r.addresses = addresses
			r.mu.Unlock()

			if err = r.UpdateState(resolver.State{Addresses: addresses}); err != nil {
				log.Errorf("failed to update resolver state: %v", err)
			}
		}
	}
}

func (r *Resolver) UpdateState(state resolver.State) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.cc.UpdateState(state)
}

func (r *Resolver) refreshAddresses() {
	services, _, err := r.consul.Health().Service(r.target.Endpoint(), "", true, nil)
	if err != nil {
		log.Errorf("failed to load healty services from consul with endpoint: %s error: %v", r.target.Endpoint(), err)
	}

	var addresses = make([]resolver.Address, len(services))
	for i, service := range services {
		addresses[i] = resolver.Address{
			Addr: fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port),
		}
	}

	r.mu.Lock()
	r.addresses = addresses
	r.mu.Unlock()

	if err = r.UpdateState(resolver.State{Addresses: addresses}); err != nil {
		log.Errorf("failed to update resolver state: %v", err)
	}
}
