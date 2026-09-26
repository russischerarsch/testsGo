package http

import (
	"context"
	"errors"
	"net/http"
	"tgtest/domain"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Handler interface {
	CreateUser(context.Context, string, string, string, string) (string, error)
	GetBalance(context.Context, string) (string, error)
}
type HandlerStruct struct {
	handler Handler
}

func CreateHandler(serv Handler) *HandlerStruct {
	return &HandlerStruct{handler: serv}
}

func (h *HandlerStruct) GetBalance(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id is missing"})
	}
	balance, err := h.handler.GetBalance(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		case pgx.ErrNoRows:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		case context.DeadlineExceeded:
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": context.DeadlineExceeded})
			return
		case domain.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidInput})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(200, balance)
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
