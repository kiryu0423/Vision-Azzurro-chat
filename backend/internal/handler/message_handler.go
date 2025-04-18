package handler

import (
	"chat-app/internal/repository"
	"chat-app/internal/service"
	"net/http"
	"strconv"

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
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDAny.(uint)

	roomIDStr := c.Param("room_id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room_id"})
		return
	}

	// 所属ユーザーであるか検証
	if err := h.RoomService.AuthorizeUser(userID, roomID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	// クエリパラメータ
	before := c.Query("before")
	limitStr := c.DefaultQuery("limit", "30")
	limit, _ := strconv.Atoi(limitStr)

	// メッセージ取得
	messages, err := h.MessageRepo.GetMessagesBefore(roomID, before, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
