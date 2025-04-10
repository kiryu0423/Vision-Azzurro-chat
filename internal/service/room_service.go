package service

import (
	"chat-app/internal/repository"
	"errors"
)

type RoomService struct {
	Repo *repository.RoomRepository
}

func NewRoomService(repo *repository.RoomRepository) *RoomService {
	return &RoomService{Repo: repo}
}

func (s *RoomService) AuthorizeUser(userID uint, roomID string) error {
	ok, err := s.Repo.InUserInRoom(userID, roomID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("unauthorized access to room")
	}
	return nil
}

func (s *RoomService) CreateRoomIfNotExists(userAID, userBID uint) (string, error) {
    room, err := s.Repo.FindRoomByUsers(userAID, userBID)
    if err == nil && room != nil && room.ID != "" {
        return room.ID, nil
    }

    newRoom, err := s.Repo.CreateRoomWithUsers([]uint{userAID, userBID})
    if err != nil {
        return "", err
    }

    return newRoom.ID, nil
}

