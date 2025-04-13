package model

import "github.com/google/uuid"

type RoomMember struct {
    RoomID uuid.UUID `json:"room_id"`
    UserID uint   `json:"user_id"`
}
