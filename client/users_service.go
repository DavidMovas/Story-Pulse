package client

import (
	"story-pulse/contracts"
)

type userWrapper struct {
	User *contracts.User `json:"user"`
}

func newWrapper() *userWrapper {
	return &userWrapper{
		User: &contracts.User{},
	}
}

func (c *Client) GetUserByID(req *contracts.GetUserByIDRequest) (*contracts.User, error) {
	var wrapper = newWrapper()

	_, err := c.client.R().
		SetResult(&wrapper).
		Get(c.path("/v1/users/%s", req.ID))

	return wrapper.User, err
}
