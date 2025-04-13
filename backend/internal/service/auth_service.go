package service

import (
	"chat-app/internal/dto"
	"chat-app/internal/model"
	"chat-app/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) Register(req dto.RegisterRequest) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user := req.ToModel(string(hashedPassword))
    return s.Repo.Create(user)
}

func (s *AuthService) Login(req dto.LoginRequest) (*model.User, error) {
	user, err := s.Repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
