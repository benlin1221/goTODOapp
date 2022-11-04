package repos

import (
	"errors"
	"m/v2/app/models"
	"m/v2/database"

	"gorm.io/gorm"
)

// FindUser searches the user's table with the condition given
func FindUser(dest interface{}, conds ...interface{}) error {
	// works
	return database.DBConn.Model(&models.User{}).Take(dest, conds...).Error
}

// FindUserByEmail searches the user's table with the email given
func FindUserByUsername(dest interface{}, username string) error {
	return FindUser(dest, "username = ?", username)
}

func FindUserByID(dest interface{}, id uint) error {
	return FindUser(dest, "id = ?", id)
}

func GetUserByID(userID uint) (*models.User, error) {
	user := &models.User{}
	if err := database.DBConn.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser create a user entry in the user's table
func CreateUser(user *models.User) error {
	err := FindUserByUsername(nil, user.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return database.DBConn.Create(user).Error
	} else {
		return errors.New("User with that username already exists")
	}
}

func ChangeUserPassword(userID uint, newPassword string) error {
	if err := database.DBConn.Model(&models.User{}).Where("ID = ?", userID).Update("Password", newPassword).Error; err != nil {
		return err
	}
	return nil
}
