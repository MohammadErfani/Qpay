package services

import (
	"Qpay/models"
	"Qpay/utils"
	"errors"
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

func CheckIsAdmin(db *gorm.DB, userID uint) error {
	var role uint8
	err := db.Model(&models.User{}).Select("role").Where("id", userID).Scan(&role).Error
	if err != nil {
		return errors.New("error getting user")
	}
	if role != models.IsAdmin {
		return errors.New("unAuthorize")
	}
	return nil
}

func ChangeUserGatewaysStatus(db *gorm.DB, user *models.User, status string) error {
	for _, gateway := range user.Gateways {
		_, err := SetStatusGateway(db, &gateway, status)
		if err != nil {
			return err
		}
	}
	return nil

}
