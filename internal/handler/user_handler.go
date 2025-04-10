package handler

import (
	"chat-app/internal/model"
	"chat-app/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserRepository *repository.UserRepository
}

func (h *UserHandler) GetAllUser(c *gin.Context) {
	users, err := h.UserRepository.FindAll()

	if err != nil {
		c.HTML(500, "users.html", gin.H{"Users": []model.User{}})
		return
	}
	
	c.HTML(200, "users.html", gin.H{"Users": users})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")

	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
	}

	id := uint(id64)

	user, err := h.UserRepository.FindByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	c.JSON(200, user)
}
