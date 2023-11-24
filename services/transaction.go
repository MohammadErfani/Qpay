package services

import (
	"Qpay/models"
	"Qpay/utils"
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
func CancelledTransaction(db *gorm.DB, TrackingCode string) error {
	var trans models.Transaction
	err := db.Where("tracking_code=?", TrackingCode).First(&trans).Error
	if err != nil {
		return errors.New("transaction Not found")
	}
	trans.Status = models.Cancelled
	db.Save(&trans)
	return nil
}
func CreateTransaction(db *gorm.DB, TrackingCode string, PaymentAmount float64, CardYear int, CardMonth int, PhoneNumber string) (*models.Transaction, error) {
	var transaction models.Transaction
	var gateway models.Gateway
	err := db.Where("tracking_code=?", TrackingCode).First(&transaction).Error
	if err != nil {
		return &models.Transaction{}, errors.New("transaction Not found")
	}
	if err = db.Preload("Commission").First(&gateway, fmt.Sprintf("%s=?", "ID"), transaction.GatewayID).Error; err != nil {
		return nil, errors.New("gateway not found")
	}

	//اینجا متصل میشیم به ماکبانک مرکزی و تراکنش رو انجام میدیم اگه ارور نداشت
	if err := utils.Transaction(PaymentAmount, CardYear, CardMonth, PhoneNumber); err != nil {
		return nil, err
	}
	return nil, nil
}
