package handlers

import (
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

type BankAccountHandler struct {
	DB     *gorm.DB
	UserID uint
}
type BankAccountResponse struct {
	Sheba        string `json:"sheba"`
	BankName     string `json:"bank_name"`
	BankLogo     string `json:"bank_logo"`
	AccountOwner string `json:"account_owner"`
	Status       string `json:"status"`
}

func (ba *BankAccountHandler) ListAllBankAccounts(ctx echo.Context) error {

	bankAccounts, err := services.GetUserBankAccounts(ba.DB, ba.UserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "You didn't add any bank account!")
	}
	var bankAccountResponses []BankAccountResponse
	for _, ba := range bankAccounts {
		bankAccountResponses = append(bankAccountResponses, SetBankAccountResponse(ba))
	}
	return ctx.JSON(http.StatusOK, bankAccountResponses)

}
func (ba *BankAccountHandler) FindBankAccount(ctx echo.Context) error {
	var bankAccount models.BankAccount
	bankAccountID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "bank account is not correct")
	}
	bankAccount, err = services.GetSpecificBankAccount(ba.DB, ba.UserID, uint(bankAccountID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "bank account does not exist!")
	}
	bankAccountResponse := SetBankAccountResponse(bankAccount)
	return ctx.JSON(http.StatusOK, bankAccountResponse)
}

func (ba *BankAccountHandler) RegisterNewBankAccount(ctx echo.Context) error {
	var req BankAccountRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	err := utils.CheckSheba(req.Sheba)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if err := ValidateUniqueBankAccount(ba.DB, &req); err != nil {
		return ctx.JSON(http.StatusConflict, err.Error())
	}
	_, err = services.CreateBankAccount(ba.DB, ba.UserID, req.Sheba)
	if err != nil {
		if err.Error() == "UnAuthorize" {
			return ctx.JSON(http.StatusForbidden, "sheba doesn't match your credential")
		}
		return ctx.JSON(http.StatusInternalServerError, "Internal server error in create bank account")
	}
	return ctx.JSON(http.StatusOK, "You're bank account is successfully registered!")
}

func (ba *BankAccountHandler) DeleteBankAccount(ctx echo.Context) error {
	bankAccountID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "bank account is not correct")
	}
	err = services.DeleteBankAccount(ba.DB, ba.UserID, uint(bankAccountID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err.Error())
	}
	return ctx.JSON(http.StatusOK, "You're bank account is successfully deleted!")
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
