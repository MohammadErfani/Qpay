package services

import (
	"Qpay/models"
	"errors"
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
func CreateTransaction(db *gorm.DB) (*models.Transaction, error) {
	return nil, nil
}
