//go:build integration

package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tgtest/serv"
	"tgtest/serv/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_integrationTest(t *testing.T) {
	repo := mocks.NewRepoInterface(t)
	repo.On("CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("1", nil)
	service, err := serv.CreateServ(repo)
	require.NoError(t, err)
	handler := CreateHandler(service)
	router := gin.Default()
	router.POST("/users", handler.CreateUser)
	server := httptest.NewServer(router)
	defer server.Close()

	userReq := `{"name":"Alice", "email":"alice@example.com", "password":"Qwerty123!", "age":"22"}`
	resp, err := http.Post(server.URL+"/users", "application/json", strings.NewReader(userReq))
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	repo.AssertExpectations(t)
}
