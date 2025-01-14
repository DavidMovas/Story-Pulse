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

		t.Logf("RESPONSE: %v\n", res)
		t.Logf("RESPONSE: %v\n", res.User)

		require.NoError(t, err)
		require.NotEmpty(t, res.User)
		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.RefreshToken)
	})
}
