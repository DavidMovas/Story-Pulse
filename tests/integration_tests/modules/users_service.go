package modules

import (
	"story-pulse/client"
	"story-pulse/contracts"
	"story-pulse/tests/integration_tests/config"
	"story-pulse/tests/integration_tests/errors"
	"testing"
)

func UsersServiceTest(t *testing.T, client *client.Client, _ *config.TestConfig) {
	t.Run("users_service.GetUserByID: not found", func(t *testing.T) {
		req := &contracts.GetUserByIDRequest{ID: 100}
		_, err := client.GetUserByID(req)

		errors.RequireNotFoundError(t, err, "user", "id", req.ID)
	})
}
