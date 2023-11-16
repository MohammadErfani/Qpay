package services

import (
	"Qpay/models"
	"Qpay/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func SetUserAndBankForBankAccount(db *gorm.DB, bankAccount *models.BankAccount) error {

	identity, bankName := utils.GetIdentityAndBank(bankAccount.Sheba)
	if identity == "" || bankName == "" {
		return errors.New("sheba is incorrect")
	}
	user, err := GetUser(db, "identity", identity)
	if err != nil {
		return errors.New("invalid sheba, identity doesn't match by sheba")
	}
	bank, err := GetBank(db, "name", bankName)
	if err != nil {
		return errors.New("invalid sheba, bank name doesn't match by sheba")
	}
	bankAccount.BankID = bank.ID
	bankAccount.UserID = user.ID
	bankAccount.AccountOwner = user.Name
	return nil
}

func GetUserBankAccounts() {

}

func GetBankAccount(db *gorm.DB, fieldName, fieldValue string) (*models.BankAccount, error) {
	var bankAccount models.BankAccount
	err := db.First(&bankAccount, fmt.Sprintf("%s=?", fieldName), fieldValue).Error
	if err != nil {
		return nil, errors.New("bank account not found")
	}
	return &bankAccount, nil
}

func CreateBankAccount(db *gorm.DB, sheba string) (*models.BankAccount, error) {
	bankAccount := models.BankAccount{Sheba: sheba}
	err := SetUserAndBankForBankAccount(db, &bankAccount)
	if err != nil {
		return nil, err
	}
	err = db.Create(&bankAccount).Error
	if err != nil {
		return nil, err
	}
	return &bankAccount, nil
}

func DeleteBankAccount() {
}
