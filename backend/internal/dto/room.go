package dto

import "time"

type RoomWithUnread struct {
  RoomID        string    `json:"room_id"`
  DisplayName   string    `json:"display_name"`
  LastMessageAt time.Time `json:"last_message_at"`
  UnreadCount   int       `json:"unread_count"`
}
