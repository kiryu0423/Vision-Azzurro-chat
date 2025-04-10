package model

import "time"

type Message struct {
    ID        uint      `json:"id"`
    RoomID    string    `json:"room_id"`
    Sender    string    `json:"sender"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}
