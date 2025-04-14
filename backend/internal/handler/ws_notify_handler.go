package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

type NotifyWSHandler struct {
    Clients map[*websocket.Conn]bool
}

func NewNotifyWSHandler() *NotifyWSHandler {
    return &NotifyWSHandler{
    Clients: make(map[*websocket.Conn]bool),
	}
}

var upgrade = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *NotifyWSHandler) Handle(c *gin.Context) {
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
    	return
	}
	h.Clients[conn] = true

  // 維持用ループ（受信は無視）
	for {
    	_, _, err := conn.ReadMessage()
    	if err != nil {
      		delete(h.Clients, conn)
      		conn.Close()
      		break
    	}
  	}
}

// メッセージ通知用
func (h *NotifyWSHandler) Broadcast(jsonData []byte) {
	for conn := range h.Clients {
		conn.WriteMessage(websocket.TextMessage, jsonData)
	}
}
