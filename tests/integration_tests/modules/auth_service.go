package modules

import (
	"brain-wave/client"
	"brain-wave/contracts"
	"brain-wave/tests/integration_tests/config"
	"brain-wave/tests/integration_tests/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func AuthServiceTest(t *testing.T, client *client.Client, _ *config.TestConfig) {

	t.Run("auth_service.Register: email is required", func(t *testing.T) {
		req := &contracts.RegisterUserRequest{
			Username: John.Username,
			Password: John.Password,
		}
		_, err := client.RegisterUser(req)
		errors.RequireBadRequestError(t, err, "email is required")
	})

	t.Run("auth_service.Register: username is required", func(t *testing.T) {
		req := &contracts.RegisterUserRequest{
			Email:    John.Email,
			Password: John.Password,
		}
		_, err := client.RegisterUser(req)
		errors.RequireBadRequestError(t, err, "username is required")
	})

	t.Run("auth_service.Register: password is required", func(t *testing.T) {
		req := &contracts.RegisterUserRequest{
			Email:    John.Email,
			Username: John.Username,
		}
		_, err := client.RegisterUser(req)
		errors.RequireBadRequestError(t, err, "password is required")
	})

	t.Run("auth_service.Register: successes", func(t *testing.T) {
		req := &contracts.RegisterUserRequest{
			Email:    John.Email,
			Username: John.Username,
			Password: John.Password,
		}
		res, err := client.RegisterUser(req)

		require.NoError(t, err)
		require.NotEmpty(t, res.User)
		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.RefreshToken)

		require.Equal(t, John.Email, res.User.Email)
		require.Equal(t, John.Username, res.User.Username)

		John.User = res.User
		John.AccessToken = res.AccessToken
		John.RefreshToken = res.RefreshToken
	})

	t.Run("auth_service.Register: several successes", func(t *testing.T) {
		reqs := []struct {
			data  *contracts.RegisterUserRequest
			owner *TestUser
		}{
			{
				data: &contracts.RegisterUserRequest{
					Email:    Markus.Email,
					Username: Markus.Username,
					Password: Markus.Password,
				},
				owner: &Markus,
			},
			{
				data: &contracts.RegisterUserRequest{
					Email:    Tommi.Email,
					Username: Tommi.Username,
					Password: Tommi.Password,
				},
				owner: &Tommi,
			},
		}

		for _, req := range reqs {
			res, err := client.RegisterUser(req.data)

			require.NoError(t, err)
			require.NotEmpty(t, res.User)
			require.NotEmpty(t, res.AccessToken)
			require.NotEmpty(t, res.RefreshToken)

			require.Equal(t, req.data.Email, res.User.Email)
			require.Equal(t, req.data.Username, res.User.Username)

			req.owner.User = res.User
			req.owner.AccessToken = res.AccessToken
			req.owner.RefreshToken = res.RefreshToken
		}

	})

	t.Run("auth_service.Register: user with email X already exists", func(t *testing.T) {
		req := &contracts.RegisterUserRequest{
			Email:    John.Email,
			Username: John.Username,
			Password: John.Password,
		}
		_, err := client.RegisterUser(req)
		errors.RequireAlreadyExistsError(t, err, "user", "email", req.Email)
	})

	t.Run("auth_service.Login: success", func(t *testing.T) {
		req := &contracts.LoginUserRequest{
			Email:    &John.Email,
			Username: &John.Username,
			Password: &John.Password,
		}

		res, err := client.LoginUser(req)
		require.NoError(t, err)

		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.RefreshToken)

		require.Equal(t, John.Email, res.User.Email)
		require.Equal(t, John.Username, res.User.Username)

		John.User = res.User
	})

	t.Run("auth_service.Login: password required", func(t *testing.T) {
		req := &contracts.LoginUserRequest{
			Email: &John.Email,
		}
		_, err := client.LoginUser(req)
		errors.RequireBadRequestError(t, err, "password is required")
	})

	t.Run("auth_service.Login: email or username required", func(t *testing.T) {
		req := &contracts.LoginUserRequest{
			Password: &John.Password,
		}
		_, err := client.LoginUser(req)
		errors.RequireBadRequestError(t, err, "email or username required")
	})
}
