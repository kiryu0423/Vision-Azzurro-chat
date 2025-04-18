package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type NotifyWSHandler struct {
	UserClients map[uint]*websocket.Conn
	RedisClient *redis.Client
}

func NewNotifyWSHandler(redisClient *redis.Client) *NotifyWSHandler {
	return &NotifyWSHandler{
		UserClients: make(map[uint]*websocket.Conn),
		RedisClient: redisClient,
	}
}

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *NotifyWSHandler) Handle(c *gin.Context) {
	userIDAny, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDAny.(uint)

	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	h.UserClients[userID] = conn

	// Redis Subscribe 開始
	go h.subscribe(userID, conn)

	// クライアントが切断するまで受信を無視して維持
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(h.UserClients, userID)
			conn.Close()
			break
		}
	}
}

func (h *NotifyWSHandler) subscribe(userID uint, conn *websocket.Conn) {
	channel := fmt.Sprintf("user:%d", userID)
	pubsub := h.RedisClient.Subscribe(context.Background(), channel)
	ch := pubsub.Channel()

	for msg := range ch {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
			log.Println("通知送信失敗:", err)
			conn.Close()
			break
		}
	}
}
