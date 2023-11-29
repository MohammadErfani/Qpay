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

// ListAllBankAccounts godoc
// @Summary List all bank accounts for the authenticated user
// @Description Retrieve a list of all bank accounts associated with the authenticated user.
// @Tags bankaccounts
// @Accept json
// @Produce json
// @Success 200 {array} BankAccountResponse "List of bank accounts"
// @Failure 400 {object} Response "{"status": "error", "message": "You didn't add any bank account!"}"
// @Security ApiKeyAuth
// @Router /bankaccount [get]
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

// FindBankAccount godoc
// @Summary Find a bank account by ID for the authenticated user
// @Description Retrieve details of a specific bank account associated with the authenticated user.
// @Tags bankaccounts
// @Accept json
// @Produce json
// @Param id path int true "Bank Account ID"
// @Success 200 {object} BankAccountResponse "Bank account details"
// @Failure 400 {object} Response "{"status": "error", "message": "Bank account is not correct"}"
// @Failure 404 {object} Response "{"status": "error", "message": "Bank account does not exist!"}"
// @Security ApiKeyAuth
// @Router /bankaccount/{id} [get]
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

// RegisterNewBankAccount godoc
// @Summary Register a new bank account for the authenticated user
// @Description Register a new bank account for the authenticated user with the provided SHEBA number.
// @Tags bankaccounts
// @Accept json
// @Produce json
// @Param bankAccountRequest body BankAccountRequest true "Bank account registration details"
// @Success 201 {object} Response "{"status": "success", "message": "You're bank account is successfully registered!"}"
// @Failure 400 {object} Response "{"status": "error", "message": "Bind Error"}"
// @Failure 400 {object} Response "{"status": "error", "message": "Invalid SHEBA format"}"
// @Failure 409 {object} Response "{"status": "error", "message": "SHEBA already exists"}"
// @Failure 403 {object} Response "{"status": "error", "message": "SHEBA doesn't match your credentials"}"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error in create bank account"}"
// @Security ApiKeyAuth
// @Router /bankaccount [post]
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

// DeleteBankAccount godoc
// @Summary Delete a bank account by ID for the authenticated user
// @Description Delete a specific bank account associated with the authenticated user.
// @Tags bankaccounts
// @Accept json
// @Produce json
// @Param id path int true "Bank Account ID"
// @Success 200 {object} Response "{"status": "success", "message": "You're bank account is successfully deleted!"}"
// @Failure 400 {object} Response "{"status": "error", "message": "Bank account is not correct"}"
// @Failure 404 {object} Response "{"status": "error", "message": "Bank account does not exist!"}"
// @Security ApiKeyAuth
// @Router /bankaccount/{id} [delete]
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
