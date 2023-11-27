package services

import (
	"Qpay/models"
	"Qpay/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func GetSpecificTransaction(db *gorm.DB, trackingCode string) (models.Transaction, error) {
	var transaction models.Transaction
	err := db.Where("tracking_code=?", trackingCode).First(&transaction).Error
	if err != nil {
		return models.Transaction{}, errors.New("transaction Not found")
	}
	return transaction, nil
}

func CreateTransaction(db *gorm.DB, gatewayID uint, amount float64, phoneNumber string, commission float64) (models.Transaction, error) {
	gateway, err := GetGateway(db, "id", fmt.Sprintf("%v", gatewayID))
	if err != nil {
		return models.Transaction{}, err
	}
	transaction := models.Transaction{
		GatewayID:        gatewayID,
		PaymentAmount:    amount,
		PhoneNumber:      phoneNumber,
		CommissionAmount: commission,
		Status:           models.NotPaid,
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

func PaymentTransaction(db *gorm.DB, TransactionID uint, CardYear int, CardMonth int, PurchaserCard string) (models.Transaction, error) {
	var transaction models.Transaction
	var gateway models.Gateway
	err := db.Where("ID=?", TransactionID).First(&transaction).Error
	if err != nil {
		return models.Transaction{}, errors.New("transaction Not found")
	}
	if err = db.Preload("Commission").Preload("BankAccount").First(&gateway, fmt.Sprintf("%s=?", "ID"), transaction.GatewayID).Error; err != nil {
		return models.Transaction{}, errors.New("gateway not found")
	}
	// for checking time
	compare := transaction.CreatedAt.Add(10 * time.Minute)
	if time.Now().After(compare) {
		fmt.Println("kheie gozasht ke")
		// todo change status of transaction and return error
	}
	//اینجا متصل میشیم به ماکبانک مرکزی و تراکنش رو انجام میدیم اگه ارور نداشت
	TrackingCode, err := utils.Transaction(transaction.PaymentAmount, CardYear, CardMonth, transaction.PhoneNumber, PurchaserCard)
	if err != nil {
		return models.Transaction{}, err
	}

	//// باید تراکنش را برای گتوی کاربر ثبت کنیم

	err = db.Model(&transaction).Updates(models.Transaction{
		PurchaserCard:        PurchaserCard,
		CardMonth:            CardMonth,
		CardYear:             CardYear,
		TrackingCode:         TrackingCode,
		Status:               models.Paid,
		PurchaserBankAccount: utils.PurchaserBankAccount(PurchaserCard),
	}).Error
	if err != nil {
		return models.Transaction{}, err
	}
	// حالا بای خروجی را مشخص کنیم
	return transaction, nil
}

func GetUserTransactions(db *gorm.DB, userID, gatewayID uint) ([]models.Transaction, error) {

	var transactions []models.Transaction

	err := db.Joins("JOIN gateways ON transactions.gateways_id = gateways.id").
		Where("gateways.user_id = ? AND gateways.id = ?", userID, gatewayID).
		Preload("User").
		Find(&transactions).Error

	if err != nil {
		return transactions, err
	}

	if len(transactions) == 0 {
		return []models.Transaction{}, errors.New("no transactions found for the specified user and gateway")
	}

	return transactions, nil
}
func FindTransaction(db *gorm.DB, userID, transactionID uint) (models.Transaction, error) {
	var transaction models.Transaction
	err := db.Joins("JOIN users ON transactions.user_id = users.id").
		Where("transactions.id = ? AND users.id = ?", transactionID, userID).
		Preload("User").
		First(&transaction).Error

	if err != nil {
		return models.Transaction{}, errors.New("transaction not found")
	}

	return transaction, nil
}
func FilterTransaction(db *gorm.DB, UserID uint, StartDate *string, EndDate *string, MinAmount *float64, MaxAmount *float64) ([]models.Transaction, error) {
	var transactions []models.Transaction
	_, err := GetTransactionsByUserID(db, UserID)
	if err != nil {
		return nil, errors.New("this user doesn't have any transaction")
	}

	query := db.Joins("JOIN gateways ON transactions.gateways_id = gateways.id").
		Where("gateways.user_id = ?", UserID).
		Preload("Gateway.User")
	if StartDate != nil {
		parsedStartTime, err := time.Parse("2006-01-02", *StartDate)
		if err != nil {
			return nil, errors.New("start date is not in correct format")
		}
		query = query.Where("created_at >= ?", parsedStartTime)
	}
	if EndDate != nil {
		parsedEndTime, err := time.Parse("2006-01-02", *EndDate)
		if err != nil {
			return nil, errors.New("start date is not in correct format")
		}
		query = query.Where("created_at <= ?", parsedEndTime)
	}
	if MinAmount != nil {
		query = query.Where("payment_amount >= ?", MinAmount)
	}
	if MaxAmount != nil {
		query = query.Where("payment_amount <= ?", MaxAmount)
	}
	err = query.Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
func GetTransactionsByUserID(db *gorm.DB, userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := db.Joins("JOIN gateways ON transactions.gateways_id = gateways.id").
		Where("gateways.user_id = ?", userID).
		Preload("Gateway.User")

	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}
