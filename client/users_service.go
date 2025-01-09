package client

import "story-pulse/contracts"

func (c *Client) GetUserByID(req *contracts.GetUserByIDRequest) (*contracts.User, error) {
	var user *contracts.User

	_, err := c.client.R().
		SetResult(&user).
		Get(c.path("/v1/users/%d", req.ID))

	return user, err
}
