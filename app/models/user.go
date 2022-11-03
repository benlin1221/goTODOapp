package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `json:"id"`
	Username string `json:"username" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Tasks    []Task `json:"tasks"`
}

type UserDTO struct {
	Username string `json:"username" gorm:"not null"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username" validate:"omitempty,min=5,max=16,alphanum"`
	Tasks    []Task `json:"tasks"`
}
