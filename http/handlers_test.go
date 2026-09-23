package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
