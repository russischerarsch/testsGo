package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tgtest/apiclient"
	"tgtest/http/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestCreateUser_withClientTest(t *testing.T) {
	service := mocks.NewHandler(t)
	service.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("1", nil)
	handler := CreateHandler(service)
	router := gin.Default()
	router.POST("/users", handler.CreateUser)
	server := httptest.NewServer(router)
	defer server.Close()

	client := apiclient.CreateClient(server.URL)

	resp, err := client.CreateUser(&apiclient.ClientCreateUserRequest{
		Name:     "Daria",
		Email:    "daria@example.com",
		Password: "Qwerty123!",
		Age:      "22",
	})
	assert.NoError(t, err)
	assert.Equal(t, "1", resp.ID)
}
