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

//func CheckTransaction(db *gorm.DB, bankAccountID uint) (*models.User, error) {
//	user, err := GetUser(db, "id", fmt.Sprintf("%v", userID))
//	if err != nil {
//		return nil, errors.New(fmt.Sprintf("user is incorrect: %v", err.Error()))
//	}
//	for _, ba := range user.BankAccounts {
//		if bankAccountID == ba.ID {
//			return user, nil
//		}
//	}
//	return nil, errors.New("bank account is incorrect")
//}
