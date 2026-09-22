package http

import (
	"errors"
	"net/http"
	"tgtest/domain"
	"tgtest/serv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *serv.Service
}

func CreateHandler(serv *serv.Service) *Handler {
	return &Handler{service: serv}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Age      int    `json:"age"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(404, gin.H{"error": "bad request"})
		return
	}
	id, err := h.service.CreateUser(c.Request.Context(), req.Name, req.Email, req.Password, req.Age)
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
