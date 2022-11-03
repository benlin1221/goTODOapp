package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title    string `validate:"omitempty,ascii" json:"title"`
	Assignee string `json:"assignee"`
	IsDone   bool   `gorm:"default:false" json:"isDone"`
	UserID   uint   `json:"userID"`
}

type TaskDTO struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title    string `validate:"omitempty,ascii" json:"title"`
	Assignee string `json:"assignee"`
	IsDone   bool   `gorm:"default:false" json:"isDone"`
	UserID   uint   `json:"userID"`
}

type TaskResponse struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title    string `validate:"omitempty,ascii" json:"title"`
	Assignee string `json:"assignee"`
	IsDone   bool   `gorm:"default:false" json:"isDone"`
	UserID   uint   `json:"userID"`
}
