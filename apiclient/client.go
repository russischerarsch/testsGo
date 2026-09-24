package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL string
	http    *http.Client
}
type ClientCreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      string `json:"age"`
}
type ClientCreateUserResponse struct {
	ID string `json:"id"`
}

func CreateClient(baseURL string) *Client {
	return &Client{baseURL: baseURL, http: http.DefaultClient}
}

func (c *Client) CreateUser(req *ClientCreateUserRequest) (*ClientCreateUserResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Post(c.baseURL+"/users", "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	var result ClientCreateUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode id to result varuable, %w", err)
	}
	return &result, nil
}
