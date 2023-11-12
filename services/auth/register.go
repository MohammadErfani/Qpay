package auth

import (
	"Qpay/models"
	"Qpay/utils"
	"errors"
	"gorm.io/gorm"
)

func GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var dbUser models.User
	result := db.First(&dbUser, "email = ?", email)
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return &dbUser, nil
}
func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	var dbUser models.User
	result := db.First(&dbUser, "username = ?", username)
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return &dbUser, nil
}

func CreateUser(db *gorm.DB, user models.User) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	result := db.Create(&user)
	if err = result.Error; err != nil {
		return nil, err
	}

	return &user, nil
}
