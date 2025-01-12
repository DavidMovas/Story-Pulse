package integration_tests

import (
	"story-pulse/client"
	"story-pulse/contracts"
	"story-pulse/tests/integration_tests/config"
	"testing"
)

func usersServiceTest(t *testing.T, client *client.Client, _ *config.TestConfig) {
	t.Run("users_service.GetUserByID: not found", func(t *testing.T) {
		req := &contracts.GetUserByIDRequest{ID: 100}
		_, err := client.GetUserByID(req)

		t.Logf("ERROR: %v\n", err)

		requireNotFoundError(t, err, "user", "id", req.ID)
	})
}
