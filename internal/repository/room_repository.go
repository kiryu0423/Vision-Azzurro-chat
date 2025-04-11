package repository

import (
	"chat-app/internal/model"

	"gorm.io/gorm"
)

type RoomRepository struct {
    DB *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
    return &RoomRepository{DB: db}
}

type RoomListItem struct {
    RoomID        string `json:"room_id"`
    OtherUserName string `json:"other_user_name"`
}

func (r *RoomRepository) InUserInRoom(userID uint, roomID string) (bool, error) {
	var count int64
	err := r.DB.
		Table("room_members").
		Where("user_id = ? AND room_id = ?", userID, roomID).
		Count(&count).Error
	return count > 0, err
}

func (r *RoomRepository) FindRoomByUsers(userAID, userBID uint) (*model.Room, error) {
    var room model.Room
    err := r.DB.Raw(`
        SELECT r.*
        FROM rooms r
        JOIN room_members rm1 ON r.id = rm1.room_id AND rm1.user_id = ?
        JOIN room_members rm2 ON r.id = rm2.room_id AND rm2.user_id = ?
        WHERE r.is_group = false
        LIMIT 1
    `, userAID, userBID).Scan(&room).Error

    if err != nil || room.ID == "" {
        return nil, err
    }
    return &room, nil
}

func (r *RoomRepository) CreateRoomWithUsers(userIDs []uint) (*model.Room, error) {
    room := &model.Room{}
    if err := r.DB.Raw(`SELECT gen_random_uuid()`).Scan(&room.ID).Error; err != nil {
        return nil, err
    }

    if err := r.DB.Exec(`
        INSERT INTO rooms (id, is_group) 
        VALUES (?, false) 
        ON CONFLICT (id) DO NOTHING
    `, room.ID).Error; err != nil {
        return nil, err
    }

    for _, uid := range userIDs {
        if err := r.DB.Exec(`
            INSERT INTO room_members (room_id, user_id) 
            VALUES (?, ?) 
            ON CONFLICT DO NOTHING
        `, room.ID, uid).Error; err != nil {
            return nil, err
        }
    }

    return room, nil
}

func (r *RoomRepository) GetRoomByUser(userID uint) ([]RoomListItem, error) {
    var rooms []RoomListItem
    err := r.DB.Raw(`
        SELECT r.id AS room_id, u.name AS other_user_name
        FROM room_members rm
        JOIN rooms r ON rm.room_id = r.id
        JOIN room_members other_rm ON r.id = other_rm.room_id AND other_rm.user_id != ?
        JOIN users u ON u.id = other_rm.user_id
        WHERE rm.user_id = ? 
        AND r.is_group = false
        AND EXISTS (
            SELECT 1 FROM messages m WHERE m.room_id = r.id
        );
    `, userID, userID).Scan(&rooms).Error

    return rooms, err
}
