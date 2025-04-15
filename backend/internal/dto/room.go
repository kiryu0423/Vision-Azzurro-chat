package dto

import "time"

type RoomWithUnread struct {
  RoomID        string    `json:"room_id"`
  DisplayName   string    `json:"display_name"`
  IsGroup       bool      `json:"is_group"`
  LastMessage   string    `json:"last_message"`
  LastMessageAt time.Time `json:"last_message_at"`
  UnreadCount   int       `json:"unread_count"`
}
