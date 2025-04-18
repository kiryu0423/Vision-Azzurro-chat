package handler

import (
	"chat-app/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// ログイン中ユーザー以外のユーザー一覧を返す
func (h *UserHandler) ListUsers(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDAny.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	users, err := h.UserService.GetSelectableUsers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// 現在ログイン中のユーザー情報を返す
func (h *UserHandler) Me(c *gin.Context) {
	userIDAny, idExists := c.Get("user_id")
	userNameAny, nameExists := c.Get("user_name")

	if !idExists || !nameExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	userID, ok := userIDAny.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	userName, ok := userNameAny.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user name"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"user_name": userName,
	})
}
