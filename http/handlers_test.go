package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tgtest/apiclient"
	"tgtest/http/mocks"

	"github.com/gin-gonic/gin"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_handlerTest(t *testing.T) {
	handler := mocks.NewHandler(t)
	handler.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("1", nil)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	userReq := `{"name":"Alice", "email":"alice@example.com", "password":"Qwerty123!", "age":"22"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userReq))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	h := &HandlerStruct{handler}
	h.CreateUser(c)
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestCreateUser_Success(t *testing.T) {
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"http://example.com/users",
		httpmock.NewJsonResponderOrPanic(http.StatusCreated, map[string]string{"id": "123"}),
	)

	client := apiclient.CreateClientWithHTTP("http://example.com", httpClient)

	resp, err := client.CreateUser(&apiclient.ClientCreateUserRequest{
		Name:     "Daria",
		Email:    "daria@example.com",
		Password: "Qwerty123!",
		Age:      "30",
	})

	require.NoError(t, err)
	assert.Equal(t, "123", resp.ID)
}

func TestCreateUser_ServerError(t *testing.T) {
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"http://example.com/users",
		httpmock.NewStringResponder(http.StatusInternalServerError, "internal error"),
	)

	client := apiclient.CreateClientWithHTTP("http://example.com", httpClient)

	_, err := client.CreateUser(&apiclient.ClientCreateUserRequest{
		Name: "John",
	})

	require.Error(t, err)
}
