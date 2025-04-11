package service

import (
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
	Repo *repository.RoomRepository
}

func NewRoomService(repo *repository.RoomRepository) *RoomService {
	return &RoomService{Repo: repo}
}

func (s *RoomService) AuthorizeUser(userID uint, roomID uuid.UUID) error {
	ok, err := s.Repo.InUserInRoom(userID, roomID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("unauthorized access to room")
	}
	return nil
}

func (s *RoomService) CreateOneToOneRoomIfNotExists(userAID, userBID uint) (uuid.UUID, error) {
    existing, err := s.Repo.FindRoomByUsers(userAID, userBID)
    if err != nil {
        return uuid.Nil, err
    }
    if existing != nil {
        return existing.ID, nil
    }

    room := &model.Room{
        ID:        uuid.New(),
        IsGroup:   false,
        CreatedAt: time.Now(),
    }

    if err := s.Repo.CreateRoom(room, []uint{userAID, userBID}); err != nil {
        return uuid.Nil, err
    }

    return room.ID, nil
}



func (s *RoomService) CreateGroupRoomIfNotExists(creatorID uint, userIDs []uint, displayName string) (uuid.UUID, error) {
    allUserIDs := append(userIDs, creatorID)
    groupKey := GenerateGroupNameFromUserIDs(allUserIDs)

    // nameによる検索（高速で確実）
    existing, err := s.Repo.FindGroupRoomByName(groupKey)
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

    if err := s.Repo.CreateRoom(room, allUserIDs); err != nil {
        return uuid.Nil, err
    }

    return room.ID, nil
}


func (s *RoomService) GetRoomsForUser(userID uint) ([]repository.RoomListItem, error) {
    return s.Repo.GetRoomByUser(userID)
}

func GenerateGroupNameFromUserIDs(userIDs []uint) string {
    sort.Slice(userIDs, func(i, j int) bool { return userIDs[i] < userIDs[j] })
    var parts []string
    for _, id := range userIDs {
        parts = append(parts, strconv.Itoa(int(id)))
    }
    return "group_" + strings.Join(parts, "_")
}
