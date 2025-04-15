package model

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
    ID        uuid.UUID    `json:"id"`
    Name      string    `json:"name"`           // グループ名 or 空
    DisplayName  string `json:"display_name"`
    IsGroup   bool      `json:"is_group"`       // true: グループ, false: 1対1
    CreatedAt time.Time `json:"created_at"`
    LastMessage  string    `json:"last_message"` 
}
