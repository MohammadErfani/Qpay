package services

import (
	"Qpay/models"
	"Qpay/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func SetUserAndBankForBankAccount(db *gorm.DB, userID uint, bankAccount *models.BankAccount) error {

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
	if userID != user.ID {
		return errors.New("UnAuthorize")
	}
	bankAccount.BankID = bank.ID
	bankAccount.UserID = user.ID
	bankAccount.AccountOwner = user.Name
	return nil
}

// GetUserBankAccounts return all user bank accounts
// Eager
/*func GetUserBankAccounts(db *gorm.DB, userID uint) ([]models.BankAccount, error) {
	user, err := GetUser(db, "id", fmt.Sprintf("%v", userID))
	if err != nil {
		return []models.BankAccount{}, err
	}
	bankAccounts := user.BankAccounts
	if len(bankAccounts) == 0 {
		return []models.BankAccount{}, errors.New("this user doesn't have any bank account")
	}
	return bankAccounts, nil
}*/
// Normal
func GetUserBankAccounts(db *gorm.DB, userID uint) ([]models.BankAccount, error) {

	var bankAccounts []models.BankAccount
	err := db.Where("user_id=?", userID).Preload("Bank").Find(&bankAccounts).Error
	if err != nil {
		return bankAccounts, err
	}

	if len(bankAccounts) == 0 {
		return []models.BankAccount{}, errors.New("this user doesn't have any bank account")
	}
	return bankAccounts, nil
}

// GetSpecificBankAccount get only check in users bank accounts and return match
// Eager
/*func GetSpecificBankAccount(db *gorm.DB, userID, bankAccountID uint) (models.BankAccount, error) {
	var bankAccount models.BankAccount
	err := db.Where("user_id=?",userID).Where("id=?",bankAccountID).First(&bankAccount).Error
	bankAccounts, err := GetUserBankAccounts(db, userID)
	if err != nil {
		return models.BankAccount{}, err
	}
	for _, bankAccount := range bankAccounts {
		if bankAccount.ID == bankAccountID {
			return bankAccount, nil
		}
	}
	return models.BankAccount{}, errors.New("bank account not found")
}*/

// normal

func GetSpecificBankAccount(db *gorm.DB, userID, bankAccountID uint) (models.BankAccount, error) {
	var bankAccount models.BankAccount
	err := db.Where("id=? AND user_id=?", bankAccountID, userID).Preload("Bank").First(&bankAccount).Error
	if err != nil {
		return models.BankAccount{}, errors.New("bank Account Not found")
	}
	return bankAccount, nil
}
func GetBankAccount(db *gorm.DB, fieldName, fieldValue string) (*models.BankAccount, error) {
	var bankAccount models.BankAccount
	err := db.First(&bankAccount, fmt.Sprintf("%s=?", fieldName), fieldValue).Error
	if err != nil {
		return nil, errors.New("bank account not found")
	}
	return &bankAccount, nil
}

// should tested
func CreateBankAccount(db *gorm.DB, userID uint, sheba string) (*models.BankAccount, error) {
	bankAccount := models.BankAccount{Sheba: sheba}
	err := SetUserAndBankForBankAccount(db, userID, &bankAccount)
	if err != nil {
		return nil, err
	}
	err = db.Create(&bankAccount).Error
	if err != nil {
		return nil, err
	}
	return &bankAccount, nil
}

// should tested
func DeleteBankAccount(db *gorm.DB, userID, bankAccountID uint) error {
	var bankAccount models.BankAccount
	err := db.Where("id = ? AND user_id = ?", bankAccountID, userID).Preload("Gateways").First(&bankAccount).Error
	if err != nil {
		return errors.New("bank account not found")
	}
	for _, g := range bankAccount.Gateways {
		g.Status = models.StatusGatewayInActive
		db.Save(&g)
	}
	//bankAccount, err := GetSpecificBankAccount(db, userID, bankAccountID)
	//if err != nil {
	//	return err
	//}

	result := db.Delete(&bankAccount)
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}
