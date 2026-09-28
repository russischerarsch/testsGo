package suits

import (
	"encoding/json"
	"fmt"
	"testing"
	"tgtest/tests/checks"
	"tgtest/tests/clients"
	"tgtest/tests/steps"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_success(t *testing.T) {
	fmt.Println(Brokers)
	httpClient := clients.CreateHttpClient("http://localhost:8080")
	kafkaClient := clients.CreateKafkaClient(Brokers, "user-event")
	postgresClient, err := clients.CreatePostgresClient(Dsn, Ctx)
	require.NoError(t, err)
	userSteps := steps.CreateUserSteps(httpClient)
	req := clients.CreateUserRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "Qwerty123!",
		Age:      "22",
	}
	resp, err := userSteps.CreateUserSuccessfully(req.Name, req.Email, req.Password, req.Age)
	require.NoError(t, err)
	id, err := postgresClient.GetUserID(Ctx, req.Name)
	assert.NoError(t, err)
	assert.Equal(t, id, resp.ID)
	checks.CheckCreateUserResponse(t, resp, "1")
	data, _ := json.Marshal(req)
	err = kafkaClient.Producer(Ctx, data)
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		message, err := kafkaClient.Consume(Ctx)
		if err == nil {
			var got clients.CreateUserRequest
			if err := json.Unmarshal(message.Value, &got); err != nil {
				return false
			}
			require.Equal(t, req, message)
			return true
		}
		return false
	}, 10*time.Second, 100*time.Millisecond)
	t.Cleanup(func() {
		postgresClient.DeleteUser(Ctx, id)
		kafkaClient.Close()
		postgresClient.Close(Ctx)
	})
}
