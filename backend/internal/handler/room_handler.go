package handler

import (
	"chat-app/internal/service"
	"chat-app/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoomHandler struct {
	RoomService *service.RoomService
	UserService *service.UserService
}

func NewRoomHandler(roomService *service.RoomService, userService *service.UserService) *RoomHandler {
	return &RoomHandler{
		RoomService: roomService,
		UserService: userService,
	}
}

type CreateRoomRequest struct {
	UserIDs     []uint `json:"user_ids"`
	DisplayName string `json:"display_name"`
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	currentUserIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	currentUserID := currentUserIDAny.(uint)

	if len(req.UserIDs) == 1 {
		targetUserID := req.UserIDs[0]
		roomID, err := h.RoomService.CreateOneToOneRoomIfNotExists(currentUserID, targetUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create 1:1 room"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"room_id": roomID})
		return
	}

	displayName := req.DisplayName
	if displayName == "" {
		userIDs := append(req.UserIDs, currentUserID)
		names, err := h.UserService.GetUserNames(userIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user names"})
			return
		}
		displayName = util.JoinNames(names)
	}

	roomID, err := h.RoomService.CreateGroupRoomIfNotExists(currentUserID, req.UserIDs, displayName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create group room"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"room_id":      roomID,
		"display_name": displayName,
	})
}

func (h *RoomHandler) ListRooms(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDAny.(uint)

	rooms, err := h.RoomService.GetUserRoomsWithUnread(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch rooms"})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *RoomHandler) MarkRoomAsRead(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDAny.(uint)

	roomID := c.Param("room_id")
	err := h.RoomService.MarkAsRead(userID, roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *RoomHandler) UpdateRoomName(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDAny.(uint)

	roomID := c.Param("room_id")
	var req struct {
		DisplayName string `json:"display_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.RoomService.UpdateRoomName(userID, roomID, req.DisplayName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *RoomHandler) GetRoomMembers(c *gin.Context) {
	roomID := c.Param("id")

	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDAny.(uint)

	parsedUUID, err := uuid.Parse(roomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
		return
	}
	if err := h.RoomService.AuthorizeUser(userID, parsedUUID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	members, err := h.RoomService.GetMembersByRoomID(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get members"})
		return
	}

	c.JSON(http.StatusOK, members)
}

func (h *RoomHandler) LeaveRoom(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDAny.(uint)

	roomID := c.Param("room_id")
	if err := h.RoomService.LeaveRoom(roomID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to leave room"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDAny.(uint)

	roomID := c.Param("room_id")
	parsedUUID, err := uuid.Parse(roomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
		return
	}
	if err := h.RoomService.AuthorizeUser(userID, parsedUUID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.RoomService.DeleteRoom(roomID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
