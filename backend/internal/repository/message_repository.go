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
	if err := r.DB.Create(message).Error; err != nil {
		return err
	}

	return r.DB.Model(&model.Room{}).
        Where("id = ?", message.RoomID).
        Update("last_message", message.Content).Error
}

// メッセージ全件取得
func (r *MessageRepository) GetMessagesByRoom(roomID uuid.UUID) ([]model.Message, error) {
    var messages []model.Message
    err := r.DB.Where("room_id = ?", roomID).Order("created_at ASC").Find(&messages).Error
    return messages, err
}

// メッセージの指定件数取得
// repository/message_repository.go
func (r *MessageRepository) GetMessagesBefore(roomID uuid.UUID, before string, limit int) ([]model.Message, error) {
	var messages []model.Message

	query := r.DB.
		Where("room_id = ?", roomID).
		Order("created_at DESC").
		Limit(limit)

	if before != "" {
		query = query.Where("created_at < ?", before)
	}

	err := query.Find(&messages).Error
	if err != nil {
		return nil, err
	}

	// フロントが昇順表示なので、昇順に並べ直す
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
