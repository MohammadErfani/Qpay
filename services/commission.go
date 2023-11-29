package services

import (
	"Qpay/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func ListAllCommission(db *gorm.DB) ([]models.Commission, error) {
	var comms []models.Commission
	err := db.Preload("Gateways").Find(&comms).Error
	if err != nil {
		return comms, errors.New("error getting commissions")
	}
	return comms, nil
}

func ListActiveCommission(db *gorm.DB) ([]models.Commission, error) {
	var comms []models.Commission
	err := db.Where("status=?", models.CommIsActive).Find(&comms).Error
	if err != nil {
		return comms, errors.New("error getting commissions")
	}
	return comms, nil
}

func GetCommission(db *gorm.DB, fieldName, fieldValue string) (*models.Commission, error) {
	var commission models.Commission
	err := db.First(&commission, fmt.Sprintf("%s=?", fieldName), fieldValue).Error
	if err != nil {
		return nil, errors.New("commission not found")
	}
	return &commission, nil
}

func CreateCommission(db *gorm.DB, amountPerTrans, percentPerTrans float64) (*models.Commission, error) {
	comm := models.Commission{
		AmountPerTrans:     amountPerTrans,
		PercentagePerTrans: percentPerTrans,
		Status:             models.CommIsActive,
	}
	err := db.Create(&comm).Error
	if err != nil {
		return nil, err
	}
	return &comm, nil
}
