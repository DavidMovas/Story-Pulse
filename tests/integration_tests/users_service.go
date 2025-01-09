package integration_tests

import (
	"story-pulse/client"
	"story-pulse/contracts"
	"testing"
	"time"
)

func usersServiceTest(t *testing.T, client *client.Client, cfg *TestConfig) {
	t.Run("users_service.GetUserByID: not found", func(t *testing.T) {
		req := &contracts.GetUserByIDRequest{ID: 100}
		_, err := client.GetUserByID(req)
		requireNotFoundError(t, err, "user", "id", 100)
	})

	time.Sleep(time.Second * 45)
}
