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

func SetUserAndBankForGateway(db *gorm.DB, userID uint, bankAccount *models.Gateway) error {

	identity, bankName := utils.GetIdentityAndBank(bankAccount.Sheba)
	if identity == "" || bankName == "" {
		return errors.New("sheba is incorrect")
	}
	user, err := GetUser(db, "identity", identity)
	if err != nil {
		return errors.New("invalid sheba, identity doesn't match by sheba")
	}
	bank, err := GetGateway(db, "name", bankName)
	if err != nil {
		return errors.New("invalid sheba, bank name doesn't match by sheba")
	}
	if userID != user.ID {
		return errors.New("UnAuthorize")
	}
	bankAccount.BankID = bank.ID
	bankAccount.UserID = user.ID
	bankAccount.AccountOwner = user.Name
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
