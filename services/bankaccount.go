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

// GetUserBankAccounts return all user bank accounts
func GetUserBankAccounts(db *gorm.DB, userID uint) ([]models.BankAccount, error) {
	user, err := GetUser(db, "id", fmt.Sprintf("%v", userID))
	if err != nil {
		return []models.BankAccount{}, err
	}
	cards := user.BankAccounts
	if len(cards) == 0 {
		return []models.BankAccount{}, errors.New("this user doesn't have any bank account")
	}
	return cards, nil
}

// GetSpecificBankAccount get only check in users bank accounts and return match
func GetSpecificBankAccount(db *gorm.DB, userID, bankAccountID uint) (models.BankAccount, error) {
	cards, err := GetUserBankAccounts(db, userID)
	if err != nil {
		return models.BankAccount{}, err
	}
	for _, card := range cards {
		if card.ID == bankAccountID {
			return card, nil
		}
	}
	return models.BankAccount{}, errors.New("bank account not found")
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

func DeleteBankAccount(db *gorm.DB, bankAccountID uint) error {
	result := db.Delete(&models.BankAccount{}, bankAccountID)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}
