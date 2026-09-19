package http

import (
	"errors"
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
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(404, gin.H{"error": "bad request"})
		return
	}
	id, err := h.service.CreateUser(c.Request.Context(), req.Name, req.Email)
	if err != nil {
		if errors.Is(err, serv.ErrInvalidName) {
			c.JSON(404, gin.H{"error": "invalid name"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"id": id})
}
