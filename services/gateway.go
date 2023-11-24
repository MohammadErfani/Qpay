package services

import (
	"Qpay/models"
	"Qpay/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func GetUserGateways(db *gorm.DB, userID uint) ([]models.Gateway, error) {

	var gateways []models.Gateway
	err := db.Where("user_id=?", userID).Preload("User").Find(&gateways).Error
	if err != nil {
		return gateways, err
	}

	if len(gateways) == 0 {
		return []models.Gateway{}, errors.New("this user doesn't have any gateway")
	}
	return gateways, nil
}

/*func SetUserAndBankForGateway(db *gorm.DB, userID uint, gateway *models.Gateway) error {

	identity, gatewayName := utils.GetIdentityAndBank(gateway.Name)
	if identity == "" || gatewayName == "" {
		return errors.New("gateway is incorrect")
	}
	user, err := GetUser(db, "identity", identity)
	if err != nil {
		return errors.New("invalid gateway, identity doesn't match by gateway")
	}
	gat, err := GetGateway(db, "name", gatewayName)
	if err != nil {
		return errors.New("invalid gateway")
	}
	if userID != user.ID {
		return errors.New("UnAuthorize")
	}
	gateway.ID = gat.ID
	gateway.UserID = user.ID
	gateway.BankAccountID = gat.BankAccountID
	return nil
}*/

func CheckUserAndBankAccountForGateway(db *gorm.DB, userID, bankAccountID uint) (*models.User, error) {
	user, err := GetUser(db, "id", fmt.Sprintf("%v", userID))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("user is incorrect: %v", err.Error()))
	}
	for _, ba := range user.BankAccounts {
		if bankAccountID == ba.ID {
			return user, nil
		}
	}
	return nil, errors.New("bank account is incorrect")
}

func CheckPersonal(db *gorm.DB, user *models.User) bool {
	for _, g := range user.Gateways {
		if g.Type == models.PersonalTypeGateway {
			return true
		}
	}
	return false
}

func CheckCommission(db *gorm.DB, commissionID uint) error {
	com, err := GetCommission(db, "id", fmt.Sprintf("%v", commissionID))
	if err != nil {
		return err
	}
	if com.Status == models.CommIsInactive {
		return errors.New("commission is inactive")
	}
	return nil
}

func SetDefaultRoute(db *gorm.DB, user *models.User, gateway *models.Gateway) {
	username := user.Username
	if username == "" {
		username = utils.GenerateRandomString(10)
	}
	if gateway.Type == models.PersonalTypeGateway {
		gateway.Route = user.Username
		return
	}
	route := fmt.Sprintf("%s_%s", username, utils.GenerateRandomString(10))
	// check if random string exist in database
	for _, err := GetGateway(db, "route", route); err == nil; {
		route = fmt.Sprintf("%s_%s", username, utils.GenerateRandomString(10))
	}
	gateway.Route = route
}

func GetSpecificGateway(db *gorm.DB, userID, gatewayID uint) (models.Gateway, error) {
	var gateway models.Gateway
	err := db.Where("id=? AND user_id=?", gatewayID, userID).Preload("User").First(&gateway).Error
	if err != nil {
		return models.Gateway{}, errors.New("gateway Not found")
	}
	return gateway, nil
}
func GetGateway(db *gorm.DB, fieldName, fieldValue string) (*models.Gateway, error) {
	var gateway models.Gateway
	err := db.First(&gateway, fmt.Sprintf("%s=?", fieldName), fieldValue).Error
	if err != nil {
		return nil, errors.New("gateway not found")
	}
	return &gateway, nil
}

func CreateGateway(db *gorm.DB, userID uint, name, logo string, bankAccountID, commissionID uint, isPersonal bool) (*models.Gateway, error) {
	// check valid user and bank account
	user, err := CheckUserAndBankAccountForGateway(db, userID, bankAccountID)
	if err != nil {
		return nil, errors.New("UnAuthorize")
	}
	// check for doesn't add multiple personal
	if isPersonal && CheckPersonal(db, user) {
		return nil, errors.New("personal error")
	}

	err = CheckCommission(db, commissionID)
	if err != nil {
		return nil, errors.New("commission error")
	}
	gateway := models.Gateway{
		Name:          name,
		Logo:          logo,
		UserID:        userID,
		BankAccountID: bankAccountID,
		CommissionID:  commissionID,
		Type:          models.SetGatewayType(isPersonal),
	}
	SetDefaultRoute(db, user, &gateway)
	err = db.Create(&gateway).Error
	if err != nil {
		return nil, err
	}
	return &gateway, nil
}
