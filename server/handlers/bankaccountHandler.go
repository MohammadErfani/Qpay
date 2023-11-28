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

type BankAccountResponse struct {
	Sheba        string `json:"sheba"`
	BankName     string `json:"bank_name"`
	BankLogo     string `json:"bank_logo"`
	AccountOwner string `json:"account_owner"`
	Status       string `json:"status"`
}

func (h *Handler) ListAllBankAccounts(ctx echo.Context) error {
	h.SetUserID(ctx)
	bankAccounts, err := services.GetUserBankAccounts(h.DB, h.UserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "You didn't add any bank account!",
		})
	}
	var bankAccountResponses []BankAccountResponse
	for _, ba := range bankAccounts {
		bankAccountResponses = append(bankAccountResponses, SetBankAccountResponse(ba))
	}
	return ctx.JSON(http.StatusOK, bankAccountResponses)

}
func (h *Handler) FindBankAccount(ctx echo.Context) error {
	h.SetUserID(ctx)
	var bankAccount models.BankAccount
	bankAccountID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "bank account is not correct",
		})
	}
	bankAccount, err = services.GetSpecificBankAccount(h.DB, h.UserID, uint(bankAccountID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "bank account does not exist!",
		})
	}
	bankAccountResponse := SetBankAccountResponse(bankAccount)
	return ctx.JSON(http.StatusOK, bankAccountResponse)
}

func (h *Handler) RegisterNewBankAccount(ctx echo.Context) error {
	h.SetUserID(ctx)
	var req BankAccountRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Bind Error",
		})
	}
	err := utils.CheckSheba(req.Sheba)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
	}
	if err = ValidateUniqueBankAccount(h.DB, &req); err != nil {
		return ctx.JSON(http.StatusConflict, Response{
			Status:  "error",
			Message: err.Error(),
		})
	}
	_, err = services.CreateBankAccount(h.DB, h.UserID, req.Sheba)
	if err != nil {
		if err.Error() == "UnAuthorize" {
			return ctx.JSON(http.StatusForbidden, Response{
				Status:  "error",
				Message: "sheba doesn't match your credential",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Internal server error in create bank account",
		})
	}
	return ctx.JSON(http.StatusCreated, Response{
		Status:  "success",
		Message: "You're bank account is successfully registered!",
	})
}

func (h *Handler) DeleteBankAccount(ctx echo.Context) error {
	h.SetUserID(ctx)
	bankAccountID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "bank account is not correct",
		})
	}
	err = services.DeleteBankAccount(h.DB, h.UserID, uint(bankAccountID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "You're bank account is successfully deleted!",
	})
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
