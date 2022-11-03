package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `json:"id"`
	Username string `json:"username" validate:"omitempty,min=5,max=16,alphanum"`
	Password string `json:"password" validate:"omitempty,min=8,max=20,alphanum"`
	Tasks    []Task `json:"tasks"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username" validate:"omitempty,min=5,max=16,alphanum"`
	Tasks    []Task `json:"tasks"`
}
