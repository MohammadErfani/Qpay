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
func CancelledTransaction(db *gorm.DB, TrackingID uint) error {
	var trans models.Transaction
	err := db.Where("ID=?", TrackingID).First(&trans).Error
	if err != nil {
		return errors.New("transaction Not found")
	}
	trans.Status = models.Cancelled
	db.Save(&trans)
	return nil
}
func CreateTransaction(db *gorm.DB, TransactionID uint, PaymentAmount float64, CardYear int, CardMonth int, PhoneNumber string, PurchaserCard string) (models.Transaction, error) {
	var transaction models.Transaction
	var gateway models.Gateway
	err := db.Where("ID=?", TransactionID).First(&transaction).Error
	if err != nil {
		return models.Transaction{}, errors.New("transaction Not found")
	}
	if err = db.Preload("Commission").Preload("BankAccount").First(&gateway, fmt.Sprintf("%s=?", "ID"), transaction.GatewayID).Error; err != nil {
		return models.Transaction{}, errors.New("gateway not found")
	}

	//اینجا متصل میشیم به ماکبانک مرکزی و تراکنش رو انجام میدیم اگه ارور نداشت
	TrackingCode, err := utils.Transaction(PaymentAmount, CardYear, CardMonth, PhoneNumber, PurchaserCard)
	if err != nil {
		return models.Transaction{}, err
	}

	// باید تراکنش را برای گتوی کاربر ثبت کنیم
	transaction = models.Transaction{
		GatewayID:            transaction.GatewayID,
		PaymentAmount:        PaymentAmount,
		CommissionAmount:     utils.ComisionCalc(PaymentAmount, transaction.Gateway.Commission.AmountPerTrans),
		PurchaserCard:        PurchaserCard,
		CardMonth:            CardMonth,
		CardYear:             CardYear,
		PhoneNumber:          PhoneNumber,
		TrackingCode:         TrackingCode,
		Status:               models.Paid,
		OwnerBankAccount:     gateway.BankAccount.AccountOwner,
		PurchaserBankAccount: utils.PurchaserBankAccount(PurchaserCard),
	}
	db.Save(&transaction)

	// حالا بای خروجی را مشخص کنیم
	return transaction, nil
}
