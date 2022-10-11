package repos

import (
	"errors"
	"m/v2/app/models"
	"m/v2/database"

	"gorm.io/gorm"
)

func GetAllTasks(pagination models.Pagination) (*models.Pagination, error) {
	db := database.DBConn
	var tasks []models.Task
	if err := db.Scopes(paginate(tasks, nil, &pagination, db)).Find(&tasks).Error; err != nil {
		return &pagination, err
	}
	pagination.Rows = tasks
	return &pagination, nil
}

func GetTaskByID(taskId uint) (models.Task, error) {
	db := database.DBConn
	var task models.Task
	err := db.Find(&task, taskId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return task, errors.New("task not found")
	}
	return task, err
}

func CreateTask(title string, assignee string) (*models.Task, error) {
	db := database.DBConn
	var tasks []models.Task
	db.Find(&tasks)
	task := &models.Task{
		ID:       uint(len(tasks)),
		Title:    title,
		Assignee: assignee,
		//TODO: find userid from assignee string
	}
	if err := db.Create(&task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func DeleteTask(taskId uint) error {
	db := database.DBConn

	var item models.Task
	if err := db.First(&item, taskId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("todo not found")
	}
	db.Where("task_id = ?", taskId).Delete(&item)
	return nil
}