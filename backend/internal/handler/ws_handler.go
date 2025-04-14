package handler

import (
	"chat-app/internal/model"
	"chat-app/internal/repository"
	"chat-app/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
	ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // 開発中は全許可、本番では制限
    },
}

var roomClients = make(map[string]map[*websocket.Conn]string)
var roomClientsMu sync.Mutex

type WebSocketHandler struct {
    MessageRepo *repository.MessageRepository
    RoomService *service.RoomService
    NotifyWSHandler  *NotifyWSHandler
}

func NewWebSocketHandler(messageRepo *repository.MessageRepository, roomService *service.RoomService, notify *NotifyWSHandler) *WebSocketHandler {
    return &WebSocketHandler{
        MessageRepo: messageRepo,
        RoomService: roomService,
        NotifyWSHandler: notify,
    }
}

func (h *WebSocketHandler) Handle(c *gin.Context) {
    session := sessions.Default(c)
    userID := session.Get("user_id")
    userName := session.Get("user_name")
    roomIDStr := c.Query("room")

    if userID == nil || roomIDStr == "" || userName == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    roomID, err := uuid.Parse(roomIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room_id"})
        return
    }

    if err := h.RoomService.AuthorizeUser(userID.(uint), roomID); err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
        return
    }

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade"})
        return
    }

    // ルームへの接続登録
    roomClientsMu.Lock()
    if roomClients[roomIDStr] == nil {
        roomClients[roomIDStr] = make(map[*websocket.Conn]string)
    }
    roomClients[roomIDStr][conn] = userName.(string)
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
            RoomID:  roomID,
            SenderID:userID.(uint),
            Sender:  userName.(string),
            Content: string(msgBytes),
        }

        if err := h.MessageRepo.SaveMessage(msg); err != nil {
            fmt.Println("DB保存失敗:", err)
        }

        // 通知用JSON構築（sender_id込み）
        notifyMsg := map[string]interface{}{
            "room_id":    msg.RoomID,
            "sender_id":  msg.SenderID,
            "sender":     msg.Sender,
            "content":    msg.Content,
            "created_at": msg.CreatedAt,
        }
        jsonNotify, _ := json.Marshal(notifyMsg)

        // ✅ 通知用WebSocketにも送る
        h.NotifyWSHandler.Broadcast(jsonNotify)

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
