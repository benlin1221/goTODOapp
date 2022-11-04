package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `json:"username" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Tasks    []Task `json:"tasks"`
}

type UserDTO struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `json:"username" gorm:"not null"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `json:"username" validate:"omitempty,min=5,max=16,alphanum"`
	Tasks    []Task `json:"tasks"`
}
