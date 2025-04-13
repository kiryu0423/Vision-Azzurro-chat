package repository

import (
	"gorm.io/gorm"
	"chat-app/internal/model"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User

	err := r.DB.First(&user, id).Error

	return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetAllExcept(userID uint) ([]model.User, error) {
	var users []model.User
	err := r.DB.Where("id != ?", userID).Find(&users).Error
	return users, err
}

func (r *UserRepository) GetUsersByIDs(ids []uint) ([]model.User, error) {
    var users []model.User
    if err := r.DB.Where("id IN ?", ids).Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}
