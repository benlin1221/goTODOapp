package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID       uint
	Title    string `validate:"omitempty,ascii"`
	Assignee string `json:"assignee"`
	IsDone   bool   `gorm:"default:false" json:"isDone"`
	UserID   uint
}

type TaskResponse struct {
	ID       uint
	Title    string `validate:"omitempty,ascii"`
	Assignee string `json:"assignee"`
	IsDone   bool   `gorm:"default:false" json:"isDone"`
	UserID   uint
}
