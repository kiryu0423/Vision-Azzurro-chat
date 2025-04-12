package handler

import (
	"chat-app/internal/service"
	"chat-app/internal/util"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
	UserIDs     []uint `json:"user_ids"`      // 招待対象（1対1なら1件）
	DisplayName string `json:"display_name"`  // グループ名（任意）
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// ログイン中のユーザーID取得（セッションやJWTから）
	session := sessions.Default(c)
    currentUserID := session.Get("user_id").(uint)

	// --- 1対1チャットの場合 ---
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

	// --- グループチャットの場合 ---
	// 表示名が空なら「ユーザーA, ユーザーB...」形式にする
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
		"room_id": roomID,
		"display_name": displayName, // ← フロントに返す
	})
}


func (h *RoomHandler) ListRooms(c *gin.Context) {
    session := sessions.Default(c)
    userID := session.Get("user_id").(uint)

    rooms, err := h.RoomService.GetRoomsForUser(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch rooms"})
        return
    }

    c.JSON(http.StatusOK, rooms)
}
