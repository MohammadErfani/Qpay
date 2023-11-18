package services

import (
	"Qpay/models"
	"Qpay/utils"
	"gorm.io/gorm"
)

func CreateAdmin(db *gorm.DB, name, username, email, password string) (*models.User, error) {

	admin := models.User{
		Name:     name,
		Username: username,
		Email:    email,
		Role:     models.IsAdmin,
	}
	hashedPassword, err := utils.HashPassword(password)
	admin.Password = hashedPassword
	if err != nil {
		return nil, err
	}
	result := db.Create(&admin)
	if err = result.Error; err != nil {
		return nil, err
	}

	return &admin, nil
}
