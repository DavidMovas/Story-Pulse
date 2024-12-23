package resolver

import "google.golang.org/grpc/resolver"

var _ resolver.Builder = (*Builder)(nil)

type Builder struct{}

func (*Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r, err := NewResolver(target, cc, opts)
	if err != nil {
		return nil, err
	}

	r.refreshAddresses()
	go r.watchUpdates()
	return r, nil
}

func (b *Builder) Scheme() string {
	return "dynamic"
}
