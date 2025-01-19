package client

import (
	"brain-wave/contracts"
)

func (c *Client) RegisterUser(req *contracts.RegisterUserRequest) (*contracts.RegisterUserResponse, error) {
	var result *contracts.RegisterUserResponse

	res, err := c.client.R().
		SetBody(req).
		SetResult(&result).
		Post(c.path("/v1/auth/register"))

	if res != nil && !res.IsError() {
		for _, cookie := range res.Cookies() {
			if cookie.Name == "refresh_token" {
				result.RefreshToken = cookie.Value
			}
		}
	}

	return result, err
}

func (c *Client) LoginUser(req *contracts.LoginUserRequest) (*contracts.LoginUserResponse, error) {
	var result *contracts.LoginUserResponse

	res, err := c.client.R().
		SetBody(req).
		SetResult(&result).
		Post(c.path("/v1/auth/login"))

	if res != nil && !res.IsError() {
		for _, cookie := range res.Cookies() {
			if cookie.Name == "refresh_token" {
				result.RefreshToken = cookie.Value
			}
		}
	}

	return result, err
}
