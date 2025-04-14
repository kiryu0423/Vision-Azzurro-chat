package model

import "time"

type RoomRead struct {
	UserID     uint      `gorm:"primaryKey"`
	RoomID     string    `gorm:"primaryKey"`
	LastReadAt time.Time `gorm:"not null"`
}
