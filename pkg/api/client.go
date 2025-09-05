package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type ApiClient interface {
	MakeRequest(context.Context, string, string, []byte) ([]byte, error)
}

type apiClient struct {
	apiUrl string
	apiKey string
}

func NewApiClient(apiUrl string, apiKey string) ApiClient {
	return &apiClient{
		apiUrl: apiUrl,
		apiKey: apiKey,
	}
}

func (c apiClient) MakeRequest(ctx context.Context, method string, url string, payload []byte) ([]byte, error) {
	request, _ := http.NewRequest(method, c.apiUrl+url, bytes.NewBuffer(payload))

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Api-Key", c.apiKey)

	client := &http.Client{}
	httpResponse, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	body, _ := io.ReadAll(httpResponse.Body)

	if httpResponse.StatusCode >= 400 && httpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("response status code: %d", httpResponse.StatusCode)
	}

	return body, nil
}
