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

// ルーム作成
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

// ルーム一覧
func (h *RoomHandler) ListRooms(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	rooms, err := h.RoomService.GetUserRoomsWithUnread(userID.(uint))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch rooms"})
		return
	}

	c.JSON(200, rooms)
}

func (h *RoomHandler) GetUserRooms(c *gin.Context) {
	session := sessions.Default(c)
    userID := session.Get("user_id").(uint)

	rooms, err := h.RoomService.GetUserRoomsWithUnread(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch rooms"})
		return
	}

	c.JSON(200, rooms)
}

// 既読管理
func (h *RoomHandler) MarkRoomAsRead(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	roomID := c.Param("room_id")

	err := h.RoomService.MarkAsRead(userID.(uint), roomID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to mark as read"})
	return
	}

	c.JSON(200, gin.H{"status": "ok"})
}

// グループ名変更
func (h *RoomHandler) UpdateRoomName(c *gin.Context) {
	userID := sessions.Default(c).Get("user_id")
	if userID == nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	roomID := c.Param("room_id")
	var req struct {
		DisplayName string `json:"display_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	err := h.RoomService.UpdateRoomName(userID.(uint), roomID, req.DisplayName)
	if err != nil {
		c.JSON(500, gin.H{"error": "update failed"})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}

// ルームメンバー取得
func (h *RoomHandler) GetRoomMembers(c *gin.Context) {
	roomID := c.Param("id")

	// ユーザーID（認証が必要なら）
	// session := sessions.Default(c)
	// userID, ok := session.Get("user_id").(uint)
	// if !ok {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }

	// 所属確認（オプション）
	// if err := h.RoomService.AuthorizeUser(userID, roomID); err != nil {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
	// 	return
	// }

	// メンバー取得
	members, err := h.RoomService.GetMembersByRoomID(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get members"})
		return
	}

	c.JSON(http.StatusOK, members)
}
