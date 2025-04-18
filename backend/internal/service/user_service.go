package service

import (
	"chat-app/internal/model"
	"chat-app/internal/repository"

	"gorm.io/gorm"
)

type UserService struct {
    DB *gorm.DB
    Repo *repository.UserRepository
}

func NewUserService(db *gorm.DB,repo *repository.UserRepository) *UserService {
    return &UserService{
        DB: db,
        Repo: repo,
    }
}

func (s *UserService) GetSelectableUsers(currentUserID uint) ([]model.User, error) {
    return s.Repo.GetAllExcept(currentUserID)
}

func (s *UserService) GetUserNames(userIDs []uint) ([]string, error) {
	var names []string
	if err := s.DB.
		Table("members").
		Select("name").
		Where("id IN ?", userIDs).
		Pluck("name", &names).Error; err != nil {
		return nil, err
	}
	return names, nil
}
