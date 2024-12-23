package resolver

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
)

const (
	consulAddress = "http://consul:8500"
	tickerTimeout = time.Second * 15
)

var _ resolver.Resolver = (*Resolver)(nil)

type Resolver struct {
	consul      *api.Client
	target      resolver.Target
	cc          resolver.ClientConn
	opts        resolver.BuildOptions
	currentAddr int
	addresses   []resolver.Address

	mu            sync.Mutex
	ticker        *time.Ticker
	tickerTimeout time.Duration
}

func (r *Resolver) ResolveNow(_ resolver.ResolveNowOptions) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.addresses) == 0 {
		r.refreshAddresses()
	}

	_ = r.updateState()
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

	log.Infof("NEW RESOLVER CLIENT FOR: %s", target.Endpoint())

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

			log.Infof("WATCHING ADDRESSES AT: %s", time.Now().String())

			var addresses = make([]resolver.Address, len(services))
			for i, service := range services {
				addresses[i] = resolver.Address{
					Addr: fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port),
				}
			}

			r.mu.Lock()
			r.addresses = addresses
			r.mu.Unlock()

			if err = r.updateState(); err != nil {
				log.Errorf("failed to update resolver state: %v", err)
			}
		}
	}
}

func (r *Resolver) updateState() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.cc.UpdateState(r.addressSelection())
}

func (r *Resolver) refreshAddresses() {
	services, _, err := r.consul.Health().Service(r.target.Endpoint(), "", true, nil)
	if err != nil {
		log.Errorf("failed to load healty services from consul with endpoint: %s error: %v", r.target.Endpoint(), err)
	}

	log.Infof("REFRESHING ADDRESSES AT: %s", time.Now().String())

	var addresses = make([]resolver.Address, len(services))
	for i, service := range services {
		addresses[i] = resolver.Address{
			Addr: fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port),
		}
	}

	r.mu.Lock()
	r.addresses = addresses
	r.mu.Unlock()

	if err = r.updateState(); err != nil {
		log.Errorf("failed to update resolver state: %v", err)
	}
}

func (r *Resolver) addressSelection() resolver.State {
	selected := r.addresses[r.currentAddr]
	r.currentAddr = (r.currentAddr + 1) % len(r.addresses)

	return resolver.State{Addresses: []resolver.Address{selected}}
}
