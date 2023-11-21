package services

import (
	"Qpay/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func GetCommission(db *gorm.DB, fieldName, fieldValue string) (*models.Commission, error) {
	var commission models.Commission
	err := db.First(&commission, fmt.Sprintf("%s=?", fieldName), fieldValue).Error
	if err != nil {
		return nil, errors.New("commission not found")
	}
	return &commission, nil
}
