package handler

import (
	"chat-app/internal/model"
	"chat-app/internal/repository"
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

// 全体チャットのための変数
var clients = make(map[*websocket.Conn]string) // conn → ユーザー名
var clientsMu sync.Mutex

type WebSocketHandler struct {
    MessageRepo *repository.MessageRepository
}

func NewWebSocketHandler(repo *repository.MessageRepository) *WebSocketHandler {
    return &WebSocketHandler{MessageRepo: repo}
}

func (h *WebSocketHandler) Handle(c *gin.Context) {
    session := sessions.Default(c)
    userID := session.Get("user_id")
    userName := session.Get("user_name")
    if userID == nil || userName == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade"})
        return
    }

	clientsMu.Lock()
	clients[conn] = userName.(string)
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
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
            RoomID:  "global",
            Sender:  userName.(string),
            Content: string(msgBytes),
        }
        if err := h.MessageRepo.SaveMessage(msg); err != nil {
            fmt.Println("DB保存失敗:", err)
        }
		

		clientsMu.Lock()
		for c := range clients {
			if err := c.WriteMessage(websocket.TextMessage, []byte(fullMessage)); err != nil {
				c.Close()
				delete(clients, c)
			}
		}

		clientsMu.Unlock()
    }
}
