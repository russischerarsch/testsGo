package steps

import "tgtest/tests/clients"

type UserSteps struct {
	client *clients.HttpClient
}

func CreateUserSteps(client *clients.HttpClient) *UserSteps {
	return &UserSteps{client: client}
}
func (c *UserSteps) CreateUserSuccessfully(name, email, password, age string) (*clients.CreateUserResponse, error) {
	req := clients.CreateUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
		Age:      age,
	}
	return clients.Post[clients.CreateUserResponse](c.client, "/users", req)
}
