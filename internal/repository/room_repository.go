package repository

import (
	"chat-app/internal/model"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type RoomRepository struct {
    DB *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
    return &RoomRepository{DB: db}
}

type RoomListItem struct {
    RoomID      uuid.UUID   `json:"room_id"`
    MemberNames pq.StringArray `json:"member_names"`  // <- 複数対応
}

func (r *RoomRepository) InUserInRoom(userID uint, roomID uuid.UUID) (bool, error) {
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

    if err != nil || room.ID == uuid.Nil {
        return nil, err
    }    
    return &room, nil
}

func (r *RoomRepository) FindGroupRoomByName(name string) (*model.Room, error) {
    var room model.Room
    err := r.DB.
        Where("is_group = ? AND name = ?", true, name).
        First(&room).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &room, err
}

func (r *RoomRepository) CreateRoom(room *model.Room, userIDs []uint) error {
    // トランザクションでまとめて処理
    return r.DB.Transaction(func(tx *gorm.DB) error {
        // rooms テーブルに挿入
        if err := tx.Create(room).Error; err != nil {
            return err
        }

        // room_members に全メンバー登録
        var members []model.RoomMember
        for _, uid := range userIDs {
            members = append(members, model.RoomMember{
                RoomID: room.ID,
                UserID: uid,
            })
        }

        if err := tx.Create(&members).Error; err != nil {
            return err
        }

        return nil
    })
}

// 作成済のルーム取得
func (r *RoomRepository) GetRoomByUser(userID uint) ([]RoomListItem, error) {
    var rooms []RoomListItem
    err := r.DB.Raw(`
        SELECT 
            r.id AS room_id,
            ARRAY_AGG(u.name ORDER BY u.name) AS member_names
        FROM rooms r
        JOIN room_members rm ON r.id = rm.room_id
        JOIN users u ON rm.user_id = u.id
        WHERE r.id IN (
            SELECT room_id FROM room_members WHERE user_id = ?
        )
        AND EXISTS (
            SELECT 1 FROM messages m WHERE m.room_id = r.id
        )
        GROUP BY r.id;
    `, userID).Scan(&rooms).Error

    return rooms, err
}
