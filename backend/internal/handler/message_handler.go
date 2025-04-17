package handler

import (
	"chat-app/internal/repository"
	"chat-app/internal/service"
	"net/http"
	"strconv"

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
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room_id"})
		return
	}

	if err := h.RoomService.AuthorizeUser(userID, roomID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	// ✅ クエリパラメータを取得
	// handler/message_handler.go
	before := c.Query("before")
	limitStr := c.DefaultQuery("limit", "30")
	limit, _ := strconv.Atoi(limitStr)

	// repository側に渡す
	messages, err := h.MessageRepo.GetMessagesBefore(roomID, before, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
