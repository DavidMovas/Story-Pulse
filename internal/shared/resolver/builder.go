package resolver

import "google.golang.org/grpc/resolver"

var _ resolver.Builder = (*resolverBuilder)(nil)

func init() { resolver.Register(&resolverBuilder{}) }

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

func (b *resolverBuilder) Scheme() string {
	return "dynamic"
}
