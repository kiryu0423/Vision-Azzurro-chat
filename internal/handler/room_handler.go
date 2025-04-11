package handler

import (
	"chat-app/internal/service"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
    RoomService *service.RoomService
}

func NewRoomHandler(rs *service.RoomService) *RoomHandler {
    return &RoomHandler{RoomService: rs}
}

type CreateRoomRequest struct {
    TargetUserID uint `json:"target_user_id"`
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
    session := sessions.Default(c)
    currentUserID := session.Get("user_id").(uint)

    var req CreateRoomRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    roomID, err := h.RoomService.CreateRoomIfNotExists(currentUserID, req.TargetUserID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "room creation failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"room_id": roomID})
}

func (h *RoomHandler) ListRooms(c *gin.Context) {
    session := sessions.Default(c)
    userID := session.Get("user_id").(uint)

    rooms, err := h.RoomService.GetRoomsForUser(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch rooms"})
    }

    c.JSON(http.StatusOK, rooms)
}
