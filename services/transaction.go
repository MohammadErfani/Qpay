package services

import (
	"Qpay/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func GetSpecificTransaction(db *gorm.DB, trackingCode string) (models.Transaction, error) {
	var transaction models.Transaction
	err := db.Where("tracking_code=?", trackingCode).First(&transaction).Error
	if err != nil {
		return models.Transaction{}, errors.New("transaction Not found")
	}
	return transaction, nil
}
func CreateTransaction(db *gorm.DB, gatewayid uint, amount float64, phonenumber string, commission float64) (models.Transaction, error) {
	gateway, err := GetGateway(db, "id", fmt.Sprintf("%v", gatewayid))
	if err != nil {
		return models.Transaction{}, err
	}
	transaction := models.Transaction{
		GatewayID:        gatewayid,
		PaymentAmount:    amount,
		PhoneNumber:      phonenumber,
		CommissionAmount: commission,
		Status:           uint8(1),
		OwnerBankAccount: gateway.BankAccount.Sheba,
	}

	err = db.Create(&transaction).Error
	if err != nil {
		return models.Transaction{}, err
	}
	return transaction, nil
}
func GetTransactionByID(db *gorm.DB, id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := db.Preload("Gateway").First(&transaction, "id=?", id).Error
	if err != nil {
		return nil, errors.New("transaction not found")
	}
	return &transaction, nil
}
