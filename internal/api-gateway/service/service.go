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

func (u *Service) Handle(serviceUrl, method, path, query string, header http.Header, body io.ReadCloser) (*http.Response, error) {
	client := http.Client{}
	requestURL := fmt.Sprintf("%s/%s", serviceUrl, path)
	if query != "" {
		requestURL += "?" + query
	}

	defer func() {
		_ = body.Close()
	}()

	request, err := http.NewRequest(method, requestURL, body)
	request.Header = header
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
