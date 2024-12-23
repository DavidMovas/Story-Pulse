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

	if token := md.Get("token")[0]; token != "" {
		metaMap["token"] = token
	}

	if userId := md.Get("userId")[0]; userId != "" {
		metaMap["userId"] = userId
	}

	return metaMap, nil
}

func (a *AuthenticateCredentials) RequireTransportSecurity() bool {
	return false
}
