package services

import (
	"Qpay/models"
	"Qpay/utils"
	"errors"
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
	// bankAccount.UserID = user.ID
	bankAccount.AccountOwner = user.Name
	return nil
}
