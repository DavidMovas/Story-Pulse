package service

import (
	"fmt"
	"io"
	"net/http"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (u *Service) Handle(serviceUrl, method, path, query string, body io.Reader) (*http.Response, error) {
	client := http.Client{}
	requestURL := fmt.Sprintf("%s/%s", serviceUrl, path)
	if query != "" {
		requestURL += "?" + query
	}

	request, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	return response, nil
}
