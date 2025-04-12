package handler

import (
	"chat-app/internal/service"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{UserService: userService}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
    session := sessions.Default(c)
    idRaw := session.Get("user_id")
    userID, ok := idRaw.(uint)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    users, err := h.UserService.GetSelectableUsers(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
        return
    }

    c.JSON(http.StatusOK, users)
}

