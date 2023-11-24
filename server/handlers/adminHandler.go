package handlers

import (
	"Qpay/models"
	"Qpay/services"
	"Qpay/utils"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type AdminRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminHandler struct {
	DB *gorm.DB
}

type CommissionRequest struct {
	AmountPerTrans  float64 `json:"amount_per_transaction"`
	PercentPerTrans float64 `json:"Percent_per_transaction"`
}

type CommissionResponse struct {
	ID              uint    `json:"id"`
	AmountPerTrans  float64 `json:"amount_per_transaction"`
	PercentPerTrans float64 `json:"Percent_per_transaction"`
	Status          string  `json:"status"`
}

func (aH *AdminHandler) AdminCreate(ctx echo.Context) error {
	var req AdminRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	if err := ValidateAdmin(&req); err != nil {
		return ctx.JSON(http.StatusForbidden, err.Error())
	}
	if err := ValidateUniqueAdmin(aH.DB, &req); err != nil {
		return ctx.JSON(http.StatusConflict, err.Error())
	}
	_, err := services.CreateAdmin(
		aH.DB,
		req.Name,
		req.Username,
		req.Email,
		req.Password,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server error in create admin")
	}
	return ctx.JSON(http.StatusCreated, RegisterResponse{Status: "success", Message: "Admin created successfully"})

}

// commission handlers

// AdminListAllCommission  for commission manager return all commission
func (aH *AdminHandler) AdminListAllCommission(ctx echo.Context) error {
	commissions, err := services.ListAllCommission(aH.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server error in getting commissions")
	}
	var commResponses []CommissionResponse
	for _, comm := range commissions {
		commResponses = append(commResponses, SetCommissionsResponse(comm))
	}
	return ctx.JSON(http.StatusOK, commResponses)

}

func (aH *AdminHandler) AdminCreateCommission(ctx echo.Context) error {
	var req CommissionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	comm, err := services.CreateCommission(aH.DB, req.AmountPerTrans, req.PercentPerTrans)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server error in creating commission")
	}
	return ctx.JSON(http.StatusCreated, SetCommissionsResponse(*comm))
}

func (aH *AdminHandler) AdminGetCommission(ctx echo.Context) error {
	commID := ctx.Param("id")
	comm, err := services.GetCommission(aH.DB, "id", fmt.Sprintf("%v", commID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "commission not found!")
	}
	return ctx.JSON(http.StatusOK, SetCommissionsResponse(*comm))
}

func (aH *AdminHandler) AdminUpdateCommission(ctx echo.Context) error {
	//TODO change amount,percent,active/inactive
	return nil
}

func (aH *AdminHandler) AdminDeleteCommission(ctx echo.Context) error {
	//TODO delete commission
	return nil
}

// user handlers

func (aH *AdminHandler) AdminListUsers(ctx echo.Context) error {
	//TODO
	return nil
}

func (aH *AdminHandler) AdminGetUser(ctx echo.Context) error {
	//TODO
	return nil
}

func (aH *AdminHandler) AdminUpdateUser(ctx echo.Context) error {
	//TODO for block unblock all user gateways,...
	return nil
}

// gateway handlers

func (aH *AdminHandler) AdminListAllGateways(ctx echo.Context) error {
	//TODO without authorize
	return nil
}

func (aH *AdminHandler) AdminGetGateway(ctx echo.Context) error {
	//TODO without authorize
	return nil
}

func (aH *AdminHandler) AdminUpdateGateway(ctx echo.Context) error {
	//TODO block unblock,...
	return nil
}

// transaction handlers

func (aH *AdminHandler) AdminListTransactions(ctx echo.Context) error {
	//TODO
	return nil
}

func (aH *AdminHandler) AdminGetTransaction(ctx echo.Context) error {
	//TODO
	return nil
}

func (aH *AdminHandler) AdminUpdateTransaction(ctx echo.Context) error {
	//TODO
	return nil
}

//validations

func ValidateAdmin(admin *AdminRequest) error {
	requiredFields := map[string]string{
		"name":     admin.Name,
		"email":    admin.Email,
		"password": admin.Password,
		"username": admin.Username,
	}
	if err := utils.IsRequired(requiredFields); err != nil {
		return err
	}
	if err := utils.IsValidEmail(admin.Email); err != nil {
		return err
	}
	return nil
}

func ValidateUniqueAdmin(db *gorm.DB, admin *AdminRequest) error {
	uniqueFields := map[string]string{
		"username": admin.Username,
		"email":    admin.Email,
	}
	for fieldName, fieldValue := range uniqueFields {
		if _, err := services.GetUser(db, fieldName, fieldValue); err == nil {
			return errors.New(fmt.Sprintf("%s already exist", fieldName))
		}
	}
	return nil
}

func SetCommissionsResponse(comm models.Commission) CommissionResponse {
	var status string
	switch comm.Status {
	case models.CommIsActive:
		status = "active"
	default:
		status = "inactive"
	}
	return CommissionResponse{
		ID:              comm.ID,
		AmountPerTrans:  comm.AmountPerTrans,
		PercentPerTrans: comm.PercentagePerTrans,
		Status:          status,
	}
}
