package handler

import (
	"chat-app/internal/model"
	"chat-app/internal/repository"
	"chat-app/internal/service"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
}

func NewWebSocketHandler(messageRepo *repository.MessageRepository, roomService *service.RoomService) *WebSocketHandler {
    return &WebSocketHandler{
        MessageRepo: messageRepo,
        RoomService: roomService,
    }
}

func (h *WebSocketHandler) Handle(c *gin.Context) {
    session := sessions.Default(c)
    userID := session.Get("user_id")
    roomID := c.Query("room")
    userName := session.Get("user_name")

    if err := h.RoomService.AuthorizeUser(userID.(uint), roomID); err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
        return
    }

    if userID == nil || roomID == "" || userName == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade"})
        return
    }

	roomClientsMu.Lock()
    if roomClients[roomID] == nil {
        roomClients[roomID] = make(map[*websocket.Conn]string)
    }
    roomClients[roomID][conn] = userName.(string)
    roomClientsMu.Unlock()

	defer func() {
        roomClientsMu.Lock()
        delete(roomClients[roomID], conn)
        if len(roomClients[roomID]) == 0 {
            delete(roomClients, roomID)
        }
        roomClientsMu.Unlock()
        conn.Close()
    }()

    for {
        _, msgBytes, err := conn.ReadMessage()
        if err != nil {
            break
        }

        // "ユーザー名: メッセージ" 形式で送信
        fullMessage := fmt.Sprintf("%s: %s", userName, string(msgBytes))

		msg := &model.Message{
            RoomID:  roomID,
            Sender:  userName.(string),
            Content: string(msgBytes),
        }
        if err := h.MessageRepo.SaveMessage(msg); err != nil {
            fmt.Println("DB保存失敗:", err)
        }
		

		roomClientsMu.Lock()
        for c := range roomClients[roomID] {
            if err := c.WriteMessage(websocket.TextMessage, []byte(fullMessage)); err != nil {
                c.Close()
                delete(roomClients[roomID], c)
            }
        }
        roomClientsMu.Unlock()
    }
}
