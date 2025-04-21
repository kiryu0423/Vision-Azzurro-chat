package handler

import (
	"chat-app/internal/model"
	"chat-app/internal/notify"
	"chat-app/internal/repository"
	"chat-app/internal/service"
	"chat-app/internal/util"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var roomClients = make(map[string]map[*websocket.Conn]string)
var roomClientsMu sync.Mutex

type WebSocketHandler struct {
	MessageRepo     *repository.MessageRepository
	RoomService     *service.RoomService
	NotifyWSHandler *NotifyWSHandler
	RedisClient     *redis.Client
}

func NewWebSocketHandler(messageRepo *repository.MessageRepository, roomService *service.RoomService, notify *NotifyWSHandler, redisClient *redis.Client) *WebSocketHandler {
	return &WebSocketHandler{
		MessageRepo:     messageRepo,
		RoomService:     roomService,
		NotifyWSHandler: notify,
		RedisClient:     redisClient,
	}
}

func (h *WebSocketHandler) Handle(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	userID, userName, err := util.ValidateJWTAndExtract(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	roomIDStr := c.Query("room")
	if roomIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing room_id"})
		return
	}

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room_id"})
		return
	}

	if err := h.RoomService.AuthorizeUser(userID, roomID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade to websocket"})
		return
	}

	roomClientsMu.Lock()
	if roomClients[roomIDStr] == nil {
		roomClients[roomIDStr] = make(map[*websocket.Conn]string)
	}
	roomClients[roomIDStr][conn] = userName
	roomClientsMu.Unlock()

	defer func() {
		roomClientsMu.Lock()
		delete(roomClients[roomIDStr], conn)
		if len(roomClients[roomIDStr]) == 0 {
			delete(roomClients, roomIDStr)
		}
		roomClientsMu.Unlock()
		conn.Close()
	}()

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			break
		}

		msg := &model.Message{
			RoomID:   roomID,
			SenderID: userID,
			Sender:   userName,
			Content:  string(msgBytes),
		}

		loc, _ := time.LoadLocation("Asia/Tokyo")
		msg.CreatedAt = time.Now().In(loc)

		if err := h.MessageRepo.SaveMessage(msg); err != nil {
			fmt.Println("DB保存失敗:", err)
		}

		// 🔔 通知送信（送信者も含めて全員）
		members, err := h.RoomService.GetMembersByRoomID(roomID.String())
		if err == nil {
			for _, m := range members {
				notifyMsg := map[string]interface{}{
					"room_id":    msg.RoomID,
					"sender_id":  msg.SenderID,
					"sender":     msg.Sender,
					"content":    msg.Content,
					"last_message": msg.Content,
					"created_at": msg.CreatedAt.Format(time.RFC3339),
					"from_self":  m.ID == msg.SenderID,
				}
				notify.PublishToUser(h.RedisClient, m.ID, notifyMsg)
			}
		}

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("メッセージのJSON変換に失敗:", err)
			continue
		}

		roomClientsMu.Lock()
		for c := range roomClients[roomIDStr] {
			if err := c.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
				c.Close()
				delete(roomClients[roomIDStr], c)
			}
		}
		roomClientsMu.Unlock()
	}
}
