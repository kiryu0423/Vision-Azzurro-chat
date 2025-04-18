package model

import "time"

type User struct {
	ID       uint   `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"-"`
    CreatedAt time.Time `json:"-"`
    UpdatedAt time.Time `json:"-"`
}


func (User) TableName() string {
    return "members"
}
