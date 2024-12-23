package options

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials"
)

var _ credentials.PerRPCCredentials = (*AuthenticateCredentials)(nil)

type AuthenticateCredentials struct {
	token string
}

func NewAuthenticateCredentials() *AuthenticateCredentials {
	return &AuthenticateCredentials{}
}

func (a *AuthenticateCredentials) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", a.token),
	}, nil
}

func (a *AuthenticateCredentials) RequireTransportSecurity() bool {
	return false
}
