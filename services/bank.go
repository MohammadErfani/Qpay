package services

import (
	"Qpay/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func GetBank(db *gorm.DB, fieldName, fieldValue string) (*models.Bank, error) {
	var bank models.Bank
	err := db.First(&bank, fmt.Sprintf("%s=?", fieldName), fieldValue).Error
	if err != nil {
		return nil, errors.New("bank not found")
	}
	return &bank, nil
}
