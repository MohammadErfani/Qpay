package handlers

import (
	"Qpay/database"
	"Qpay/models"
	"Qpay/services"
	"Qpay/utils"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type BankAccountRequest struct {
	Sheba string `json:"sheba" xml:"sheba" form:"sheba" query:"sheba"`
}

type BankAccountResponse struct {
	Sheba        string `json:"sheba"`
	BankName     string `json:"bank_name"`
	BankLogo     string `json:"bank_logo"`
	AccountOwner string `json:"account_owner"`
	Status       string `json:"status"`
}

func ListAllCards(ctx echo.Context) error {
	db := database.DB()
	var userID uint = 1
	bankAccounts, err := services.GetUserBankAccounts(db, userID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "You Aren't Add Any Card. Please add a card!")
	}
	var bankAccountResponses []BankAccountResponse
	for _, ba := range bankAccounts {
		bankAccountResponses = append(bankAccountResponses, SetBankAccountResponse(ba))
	}
	//return json.NewEncoder(ctx.Response()).Encode(&bankAccountResponses)
	return ctx.JSON(http.StatusOK, bankAccountResponses)

}
func FindCard(ctx echo.Context) error {
	db := database.DB()
	var bankAccount models.BankAccount
	var userID uint = 1
	bankAccountID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "bank account is not correct")
	}
	bankAccount, err = services.GetSpecificBankAccount(db, userID, uint(bankAccountID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "You're Card is not exist!")
	}
	bankAccountResponse := SetBankAccountResponse(bankAccount)
	//return json.NewEncoder(ctx.Response()).Encode(&bankAccountResponse)
	return ctx.JSON(http.StatusOK, bankAccountResponse)
}

func RegisterNewCard(ctx echo.Context) error {
	db := database.DB()
	var req BankAccountRequest
	//var userID uint = 1 //user id ro bayad tashkhis bedim o inja vared konim.
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	err := utils.CheckSheba(req.Sheba)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if err := ValidateUniqueBankAccount(db, &req); err != nil {
		return ctx.JSON(http.StatusConflict, err.Error())
	}
	_, err = services.CreateBankAccount(db, req.Sheba)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server error in create bank account")
	}
	return ctx.JSON(http.StatusOK, "You're card is successfully registered!")
}

func DeleteCard(ctx echo.Context) error {
	db := database.DB()
	cardID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "bank account is not correct")
	}
	err = services.DeleteBankAccount(db, uint(cardID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err.Error())
	}
	return ctx.JSON(http.StatusOK, "You're Card is successfully deleted!")
}

func ValidateUniqueBankAccount(db *gorm.DB, bankAccount *BankAccountRequest) error {
	if _, err := services.GetBankAccount(db, "sheba", bankAccount.Sheba); err == nil {
		return errors.New("sheba already exist")
	}
	return nil
}

func SetBankAccountResponse(bankAccount models.BankAccount) BankAccountResponse {
	var status string
	if bankAccount.Status == models.StatusBankAccountActive {
		status = "active"
	} else {
		status = "disable"
	}
	return BankAccountResponse{
		Sheba:        bankAccount.Sheba,
		BankName:     bankAccount.Bank.Name,
		BankLogo:     bankAccount.Bank.Logo,
		AccountOwner: bankAccount.AccountOwner,
		Status:       status,
	}
}
