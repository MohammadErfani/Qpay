package services

import (
	"Qpay/models"
	"Qpay/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func GetUserGateways(db *gorm.DB, userID uint) ([]models.Gateway, error) {

	var gateways []models.Gateway
	err := db.Where("user_id=?", userID).Preload("User").Find(&gateways).Error
	if err != nil {
		return gateways, err
	}

	if len(gateways) == 0 {
		return []models.Gateway{}, errors.New("this user doesn't have any gateway")
	}
	return gateways, nil
}

func SetUserAndBankForGateway(db *gorm.DB, userID uint, gateway *models.Gateway) error {

	identity, gatewayName := utils.GetIdentityAndBank(gateway.Name)
	if identity == "" || gatewayName == "" {
		return errors.New("gateway is incorrect")
	}
	user, err := GetUser(db, "identity", identity)
	if err != nil {
		return errors.New("invalid gateway, identity doesn't match by gateway")
	}
	gat, err := GetGateway(db, "name", gatewayName)
	if err != nil {
		return errors.New("invalid gateway")
	}
	if userID != user.ID {
		return errors.New("UnAuthorize")
	}
	gateway.ID = gat.ID
	gateway.UserID = user.ID
	gateway.BankAccountID = gat.BankAccountID
	return nil
}

// normal
func GetSpecificGateway(db *gorm.DB, userID, gatewayID uint) (models.Gateway, error) {
	var gateway models.Gateway
	err := db.Where("id=? AND user_id=?", gatewayID, userID).Preload("User").First(&gateway).Error
	if err != nil {
		return models.Gateway{}, errors.New("Gateway Not found")
	}
	return gateway, nil
}
func GetGateway(db *gorm.DB, fieldName, fieldValue string) (*models.Gateway, error) {
	var gateway models.Gateway
	err := db.First(&gateway, fmt.Sprintf("%s=?", fieldName), fieldValue).Error
	if err != nil {
		return nil, errors.New("gateway not found")
	}
	return &gateway, nil
}

// should tested
func CreateGateway(db *gorm.DB, userID uint, name string) (*models.Gateway, error) {
	gateway := models.Gateway{Name: name}
	err := SetUserAndBankForGateway(db, userID, &gateway)
	if err != nil {
		return nil, err
	}
	err = db.Create(&gateway).Error
	if err != nil {
		return nil, err
	}
	return &gateway, nil
}
