package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HttpClient struct {
	client  *http.Client
	baseURL string
}
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      string `json:"age"`
}
type CreateUserResponse struct {
	ID string `json:"id"`
}

func CreateHttpClient(url string) *HttpClient {
	return &HttpClient{
		baseURL: url,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}
func Post[T any](c *HttpClient, path string, req interface{}) (*T, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body, %w", err)
	}
	resp, err := c.client.Post(c.baseURL+path, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to send http request, %w", err)
	}
	defer resp.Body.Close()
	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode responde body, %w", err)
	}
	return &result, nil
}
