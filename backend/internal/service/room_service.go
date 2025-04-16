package service

import (
	"chat-app/internal/dto"
	"chat-app/internal/model"
	"chat-app/internal/repository"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RoomService struct {
	rRepo *repository.RoomRepository
    uRepo *repository.UserRepository
}

func NewRoomService(roomRepo *repository.RoomRepository, userRepo *repository.UserRepository) *RoomService {
	return &RoomService{
        rRepo: roomRepo,
        uRepo: userRepo,
    }
}

// ルームメンバーにユーザーがいるか確認
func (s *RoomService) AuthorizeUser(userID uint, roomID uuid.UUID) error {
	ok, err := s.rRepo.InUserInRoom(userID, roomID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("unauthorized access to room")
	}
	return nil
}

func (s *RoomService) CreateOneToOneRoomIfNotExists(userAID, userBID uint) (uuid.UUID, error) {
    existing, err := s.rRepo.FindRoomByUsers(userAID, userBID)
    if err != nil {
        return uuid.Nil, err
    }
    if existing != nil {
        return existing.ID, nil
    }

    // 両ユーザー名を取得
    users, err := s.uRepo.GetUsersByIDs([]uint{userAID, userBID})
    if err != nil || len(users) < 2 {
        return uuid.Nil, errors.New("failed to get user names")
    }

    // 表示名作成（昇順で安定化）
    sort.Slice(users, func(i, j int) bool { return users[i].Name < users[j].Name })
    var nameParts []string
    for _, u := range users {
        nameParts = append(nameParts, u.Name)
    }

    displayName := strings.Join(nameParts, ", ")
    groupKey := "oneonone_" + GenerateGroupNameFromUserIDs([]uint{userAID, userBID})

    room := &model.Room{
        ID:          uuid.New(),
        IsGroup:     false,
        Name:        groupKey,      // 内部識別キー
        DisplayName: displayName,   // 表示用（A, B）
        CreatedAt:   time.Now(),
    }

    if err := s.rRepo.CreateRoom(room, []uint{userAID, userBID}); err != nil {
        return uuid.Nil, err
    }

    return room.ID, nil
}



func (s *RoomService) CreateGroupRoomIfNotExists(creatorID uint, userIDs []uint, displayName string) (uuid.UUID, error) {
    allUserIDs := append(userIDs, creatorID)
    groupKey := GenerateGroupNameFromUserIDs(allUserIDs)

    // nameによる検索（高速で確実）
    existing, err := s.rRepo.FindGroupRoomByName(groupKey)
    if err != nil {
        return uuid.Nil, err
    }
    if existing != nil {
        return existing.ID, nil
    }

    room := &model.Room{
        ID:          uuid.New(),
        IsGroup:     true,
        Name:        groupKey,
        DisplayName: displayName, // ← 初期表示用
        CreatedAt:   time.Now(),
    }

    if err := s.rRepo.CreateRoom(room, allUserIDs); err != nil {
        return uuid.Nil, err
    }

    return room.ID, nil
}


func (s *RoomService) GetRoomsForUser(userID uint) ([]repository.RoomListItem, error) {
    return s.rRepo.GetRoomByUser(userID)
}

func GenerateGroupNameFromUserIDs(userIDs []uint) string {
    sort.Slice(userIDs, func(i, j int) bool { return userIDs[i] < userIDs[j] })
    var parts []string
    for _, id := range userIDs {
        parts = append(parts, strconv.Itoa(int(id)))
    }
    return "group_" + strings.Join(parts, "_")
}

// 未読管理
func (s *RoomService) GetUserRoomsWithUnread(userID uint) ([]dto.RoomWithUnread, error) {
    rooms, err := s.rRepo.GetRoomsWithUnreadCount(userID)
    if err != nil {
        return nil, err
    }
    user, err := s.uRepo.FindByID(userID)
    if err != nil {
        return nil, err
    }

    for i, room := range rooms {

        if !room.IsGroup {
            names := strings.Split(room.DisplayName, ",")
            var others []string
            for _, name := range names {
                if strings.TrimSpace(name) != user.Name {
                    others = append(others, strings.TrimSpace(name))
                }
            }
            rooms[i].DisplayName = strings.Join(others, ", ")
        }
    }

    return rooms, nil
}


// 既読管理
func (s *RoomService) MarkAsRead(userID uint, roomID string) error {
    return s.rRepo.UpsertRoomRead(userID, roomID)
}

// グループ名変更
func (s *RoomService) UpdateRoomName(userID uint, roomID string, name string) error {
    // 権限チェックがあればここで（省略可）
    return s.rRepo.UpdateDisplayName(roomID, name)
}

// ルームメンバー取得
func (s *RoomService) GetMembersByRoomID(roomID string) ([]dto.UserSummary, error) {
	return s.rRepo.GetRoomMembers(roomID)
}

// グループ退会
func (s *RoomService) LeaveRoom(roomID string, userID uint) error {
	return s.rRepo.RemoveMember(roomID, userID)
}

// グループ削除
func (s *RoomService) DeleteRoom(roomID string) error {
	return s.rRepo.DeleteRoom(roomID)
}
