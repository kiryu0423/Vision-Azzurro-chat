package handler

import (
	"chat-app/internal/model"
	"chat-app/internal/repository"
	"chat-app/internal/service"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type MessageHandler struct {
	MessageRepo *repository.MessageRepository
	RoomService *service.RoomService
}

func NewMessageHandler(messageRepo *repository.MessageRepository, roomService *service.RoomService) *MessageHandler {
	return &MessageHandler{
        MessageRepo: messageRepo,
        RoomService: roomService,
    }
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	session := sessions.Default(c)
    userID := session.Get("user_id").(uint)
    roomIDStr := c.Param("room_id")
	roomID,_ := uuid.Parse(roomIDStr)

    if err := h.RoomService.AuthorizeUser(userID, roomIDStr); err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
        return
    }

	messages, err := h.MessageRepo.GetMessagesByRoom(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to messages"})
		return
	}

	if messages == nil {
		messages = []model.Message{}
	}

	c.JSON(http.StatusOK, messages)
}
