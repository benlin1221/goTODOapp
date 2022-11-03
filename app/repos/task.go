package repos

import (
	"errors"
	"m/v2/app/models"
	"m/v2/database"

	"gorm.io/gorm"
)

func GetAllTasks() (*[]models.Task, error) {
	db := database.DBConn
	var tasks []models.Task
	err := db.Find(&tasks).Error
	return &tasks, err
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
		//ID:       uint(len(tasks)),
		Title:    title,
		Assignee: assignee,
		//TODO: find userid from assignee string
	}
	if err := db.Create(&task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func UpdateTask(taskID uint, title string, assignee string, isDone bool) (*models.Task, error) {
	var task = models.Task{Title: title, Assignee: assignee}
	task.ID = taskID
	err := database.DBConn.Model(&task).Select("Title", "Assignee", "IsDone").Where("id = ?", taskID).Updates(models.Task{Title: title, Assignee: assignee, IsDone: isDone}).Error
	return &task, err

}

func DeleteTask(taskID uint) error {
	var item models.Task
	err := database.DBConn.Delete(&item, taskID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("task not found")
	}
	return nil
}
