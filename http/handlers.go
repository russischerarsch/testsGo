package http

import (
	"context"
	"errors"
	"net/http"
	"tgtest/domain"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	CreateUser(context.Context, string, string, string, string) (string, error)
}
type HandlerStruct struct {
	handler Handler
}

func CreateHandler(serv Handler) *HandlerStruct {
	return &HandlerStruct{handler: serv}
}

func (h *HandlerStruct) CreateUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Age      string `json:"age"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(404, gin.H{"error": "bad request"})
		return
	}
	id, err := h.handler.CreateUser(c.Request.Context(), req.Name, req.Email, req.Password, req.Age)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidName) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
