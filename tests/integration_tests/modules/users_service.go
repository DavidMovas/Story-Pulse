package modules

import (
	"github.com/stretchr/testify/require"
	"story-pulse/client"
	"story-pulse/contracts"
	"story-pulse/tests/integration_tests/config"
	"story-pulse/tests/integration_tests/errors"
	"testing"
)

const (
	DefaultPassword = "passPASS123!@"
)

type TestUser struct {
	*contracts.User
	Password     string
	AccessToken  string
	RefreshToken string
}

var (
	John = &TestUser{
		User: &contracts.User{
			Email:    "john@gmail.com",
			Username: "john",
		},
		Password: DefaultPassword,
	}

	Markus = TestUser{
		User: &contracts.User{
			Email:    "markus@gmail.com",
			Username: "markus",
		},
		Password: DefaultPassword,
	}

	Tommi = TestUser{
		User: &contracts.User{
			Email:    "tommi@gmail.com",
			Username: "tommi",
		},
		Password: DefaultPassword,
	}
)

func UsersServiceTest(t *testing.T, client *client.Client, _ *config.TestConfig) {

	t.Run("users_service.GetUserByID: not found", func(t *testing.T) {
		req := &contracts.GetUserByIDRequest{ID: "100"}
		_, err := client.GetUserByID(req)

		errors.RequireNotFoundError(t, err, "user", "id", req.ID)
	})

	t.Run("users_service.GetUserByID: successes", func(t *testing.T) {
		req := &contracts.GetUserByIDRequest{ID: John.ID}
		res, err := client.GetUserByID(req)

		require.NoError(t, err)
		require.Equal(t, John.ID, res.ID)
		require.Equal(t, John.Email, res.Email)
		require.Equal(t, John.Username, res.Username)
	})
}
