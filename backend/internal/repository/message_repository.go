package repository

import (
	"chat-app/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type MessageRepository struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

func (r *MessageRepository) SaveMessage(message *model.Message) error {
	return r.DB.Create(message).Error
}

func (r *MessageRepository) GetMessagesByRoom(roomID uuid.UUID) ([]model.Message, error) {
    var messages []model.Message
    err := r.DB.Where("room_id = ?", roomID).Order("created_at ASC").Find(&messages).Error
    return messages, err
}
