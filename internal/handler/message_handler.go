package handler

import (
	"chat-app/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)


type MessageHandler struct {
	MessageRepo *repository.MessageRepository
}

func NewMessageHandler(repo *repository.MessageRepository) *MessageHandler {
	return &MessageHandler{MessageRepo: repo}
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	roomID := c.Param("room_id")

	messages, err := h.MessageRepo.GetMessagesByRoom(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
