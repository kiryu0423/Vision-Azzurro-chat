package dto

import (
    "chat-app/internal/model"
)

type RegisterRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

func (r *RegisterRequest) ToModel(hashedPassword string) *model.User {
    return &model.User{
        Name:     r.Name,
        Email:    r.Email,
        Password: hashedPassword,
    }
}
