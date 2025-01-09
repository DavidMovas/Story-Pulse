package client

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"story-pulse/contracts"
)

type Client struct {
	client  *resty.Client
	Address string
}

func NewClient(address string) *Client {
	hc := &http.Client{}
	rc := resty.NewWithClient(hc)
	rc.OnAfterResponse(func(_ *resty.Client, response *resty.Response) error {
		if response.IsError() {
			herr := contracts.HTTPError{}
			_ = json.Unmarshal(response.Body(), &herr)

			return &Error{Code: response.StatusCode(), Message: herr.Message}
		}
		return nil
	})

	return &Client{
		client:  rc,
		Address: address,
	}
}

func (c *Client) path(f string, args ...any) string {
	return fmt.Sprintf(c.Address+f, args...)
}
