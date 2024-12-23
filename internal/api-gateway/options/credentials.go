package options

import (
	"context"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var _ credentials.PerRPCCredentials = (*AuthenticateCredentials)(nil)

type AuthenticateCredentials struct{}

func NewAuthenticateCredentials() *AuthenticateCredentials {
	return &AuthenticateCredentials{}
}

func (a *AuthenticateCredentials) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	var metaMap = make(map[string]string)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return metaMap, nil
	}

	if tokens := md.Get("token"); len(tokens) > 0 {
		metaMap["token"] = tokens[0]
	}

	if userIds := md.Get("userId"); len(userIds) > 0 {
		metaMap["userId"] = userIds[0]
	}

	return metaMap, nil
}

func (a *AuthenticateCredentials) RequireTransportSecurity() bool {
	return false
}
