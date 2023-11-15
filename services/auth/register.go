package auth

import (
	"Qpay/models"
	"Qpay/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func GetUser(db *gorm.DB, fieldName, fieldValue string) (*models.User, error) {
	var user models.User
	err := db.First(&user, fmt.Sprintf("%s=?", fieldName), fieldValue).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
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
func GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var dbUser models.User
	err := db.First(&dbUser, "email=?", email).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &dbUser, nil
}
func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	var dbUser models.User
	err := db.First(&dbUser, "username=?", username).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &dbUser, nil
}
