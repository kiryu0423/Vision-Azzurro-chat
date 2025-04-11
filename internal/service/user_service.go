package service

import (
    "chat-app/internal/model"
    "chat-app/internal/repository"
)

type UserService struct {
    Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
    return &UserService{Repo: repo}
}

func (s *UserService) GetSelectableUsers(currentUserID uint) ([]model.User, error) {
    return s.Repo.GetAllExcept(currentUserID)
}
