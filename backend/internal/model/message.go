package model

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
    ID        uint      `json:"id"`
    RoomID    uuid.UUID `json:"room_id"`
    SenderID  uint      `json:"sender_id"` 
    Sender    string    `json:"sender"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}
