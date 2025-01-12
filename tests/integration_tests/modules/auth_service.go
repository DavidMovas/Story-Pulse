package modules

import (
	"github.com/stretchr/testify/require"
	"story-pulse/client"
	"story-pulse/contracts"
	"story-pulse/tests/integration_tests/config"
	"testing"
)

func AuthServiceTest(t *testing.T, client *client.Client, _ *config.TestConfig) {
	t.Run("auth_service.Register: successes", func(t *testing.T) {
		req := &contracts.RegisterUserRequest{
			Email:    "test_1@test.com",
			Username: "test_1",
			Password: "testPASS123!@",
		}
		res, err := client.RegisterUser(req)
		require.NoError(t, err)
		require.NotEmpty(t, res.User)
		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.RefreshToken)

		t.Logf("RESPONSE: %v\n", res)
	})
}

/*
panic: runtime error: index out of range [0] with length 0
2025-01-12 22:50:50
2025-01-12 22:50:50 goroutine 23 [running]:
2025-01-12 22:50:50 story-pulse/internal/shared/resolver.(*Resolver).addressSelection(...)
2025-01-12 22:50:50     /app/internal/shared/resolver/resolver.go:131
2025-01-12 22:50:50 story-pulse/internal/shared/resolver.(*Resolver).updateState(0xc000342900)
2025-01-12 22:50:50     /app/internal/shared/resolver/resolver.go:101 +0x2cd
2025-01-12 22:50:50 story-pulse/internal/shared/resolver.(*Resolver).refreshAddresses(0xc000342900)
2025-01-12 22:50:50     /app/internal/shared/resolver/resolver.go:121 +0x445
2025-01-12 22:50:50 story-pulse/internal/shared/resolver.(*Builder).Build(0x7fef8856c108?, {{{0xc0001ba570, 0x7}, {0x0, 0x0}, 0x0, {0xc0001ba57a, 0x0}, {0xc0001ba57a, 0xd}, ...}}, ...)
2025-01-12 22:50:50     /app/internal/shared/resolver/builder.go:15 +0xb2
2025-01-12 22:50:50 google.golang.org/grpc.(*ccResolverWrapper).start.func1({0xc927a0?, 0xc0001ce820?})
2025-01-12 22:50:50     /go/pkg/mod/google.golang.org/grpc@v1.69.2/resolver_wrapper.go:81 +0x1dc
2025-01-12 22:50:50 google.golang.org/grpc/internal/grpcsync.(*CallbackSerializer).run(0xc000187bd0, {0xc927a0, 0xc0001ce820})
2025-01-12 22:50:50     /go/pkg/mod/google.golang.org/grpc@v1.69.2/internal/grpcsync/callback_serializer.go:94 +0x174
2025-01-12 22:50:50 created by google.golang.org/grpc/internal/grpcsync.NewCallbackSerializer in goroutine 1
2025-01-12 22:50:50     /go/pkg/mod/google.golang.org/grpc@v1.69.2/internal/grpcsync/callback_serializer.go:52 +0x11a
*/
