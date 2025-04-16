package repository

import (
	"chat-app/internal/dto"
	"chat-app/internal/model"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoomRepository struct {
    DB *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
    return &RoomRepository{DB: db}
}

type RoomListItem struct {
    RoomID      uuid.UUID `json:"room_id"`
    DisplayName string    `json:"display_name"`
    LastMessageAt time.Time `json:"last_message_at"`
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
        r.display_name,
        r.last_message,
        MAX(m.created_at) AS last_message_at
    FROM rooms r
    JOIN room_members rm ON r.id = rm.room_id
    LEFT JOIN messages m ON m.room_id = r.id
    WHERE rm.user_id = ?
    GROUP BY r.id, r.display_name
    ORDER BY last_message_at DESC NULLS LAST;
    `, userID).Scan(&rooms).Error

    return rooms, err
}

// 未読管理
func (r *RoomRepository) GetRoomsWithUnreadCount(userID uint) ([]dto.RoomWithUnread, error) {
    var result []dto.RoomWithUnread

    query := `
        SELECT
        r.id AS room_id,
        r.display_name,
        r.is_group,
        r.last_message,
        MAX(m.created_at) AS last_message_at,
        COUNT(CASE
            WHEN m.created_at > COALESCE(rr.last_read_at, '1970-01-01')
                AND m.sender_id != ? THEN 1
            ELSE NULL
        END) AS unread_count
        FROM rooms r
        JOIN room_members rm ON r.id = rm.room_id
        LEFT JOIN messages m ON m.room_id = r.id
        LEFT JOIN room_reads rr ON rr.room_id = r.id AND rr.user_id = ?
        WHERE rm.user_id = ?
        GROUP BY r.id, r.display_name, rr.last_read_at
    `

    if err := r.DB.Raw(query, userID, userID, userID).Scan(&result).Error; err != nil {
        return nil, err
    }

    return result, nil
}

// 既読管理
func (r *RoomRepository) UpsertRoomRead(userID uint, roomID string) error {
    read := model.RoomRead{
        UserID:     userID,
        RoomID:     roomID,
        LastReadAt: time.Now(),
    }

    return r.DB.
        Clauses(clause.OnConflict{
            Columns:   []clause.Column{{Name: "user_id"}, {Name: "room_id"}},
            DoUpdates: clause.AssignmentColumns([]string{"last_read_at"}),
        }).
        Create(&read).Error
}

// グループ名変更
func (r *RoomRepository) UpdateDisplayName(roomID string, name string) error {
    return r.DB.Model(&model.Room{}).
        Where("id = ? AND is_group = true", roomID).
        Update("display_name", name).Error
}

// ルームメンバー取得
func (r *RoomRepository) GetRoomMembers(roomID string) ([]dto.UserSummary, error) {
	var users []dto.UserSummary
	err := r.DB.Raw(`
		SELECT u.id, u.name
		FROM users u
		JOIN room_members rm ON rm.user_id = u.id
		WHERE rm.room_id = ?
		ORDER BY u.name
	`, roomID).Scan(&users).Error

	return users, err
}
