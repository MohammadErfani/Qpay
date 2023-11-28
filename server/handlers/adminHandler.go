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

func (h *Handler) AdminCreate(ctx echo.Context) error {
	var req AdminRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	if err := ValidateAdmin(&req); err != nil {
		return ctx.JSON(http.StatusForbidden, err.Error())
	}
	if err := ValidateUniqueAdmin(h.DB, &req); err != nil {
		return ctx.JSON(http.StatusConflict, err.Error())
	}
	_, err := services.CreateAdmin(
		h.DB,
		req.Name,
		req.Username,
		req.Email,
		req.Password,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server error in create admin")
	}
	return ctx.JSON(http.StatusCreated, Response{Status: "success", Message: "Admin created successfully"})

}

// commission handlers

// AdminListAllCommission  for commission manager return all commission
func (h *Handler) AdminListAllCommission(ctx echo.Context) error {
	commissions, err := services.ListAllCommission(h.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server error in getting commissions")
	}
	var commResponses []CommissionResponse
	for _, comm := range commissions {
		commResponses = append(commResponses, SetCommissionsResponse(comm))
	}
	return ctx.JSON(http.StatusOK, commResponses)

}

func (h *Handler) AdminCreateCommission(ctx echo.Context) error {
	var req CommissionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	comm, err := services.CreateCommission(h.DB, req.AmountPerTrans, req.PercentPerTrans)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server error in creating commission")
	}
	return ctx.JSON(http.StatusCreated, SetCommissionsResponse(*comm))
}

func (h *Handler) AdminGetCommission(ctx echo.Context) error {
	commID := ctx.Param("id")
	comm, err := services.GetCommission(h.DB, "id", fmt.Sprintf("%v", commID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "commission not found!")
	}
	return ctx.JSON(http.StatusOK, SetCommissionsResponse(*comm))
}

// user handlers

func (h *Handler) AdminListUsers(ctx echo.Context) error {
	//TODO
	return nil
}

func (h *Handler) AdminGetUser(ctx echo.Context) error {
	//TODO
	return nil
}

func (h *Handler) AdminUpdateUser(ctx echo.Context) error {
	//TODO for block unblock all user gateways,...
	return nil
}

// gateway handlers

func (h *Handler) AdminListAllGateways(ctx echo.Context) error {
	gateways, err := services.ListAllGateway(h.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server error in getting gateways")
	}
	var gatewayResponses []GatewayResponse
	for _, gateway := range gateways {
		gatewayResponses = append(gatewayResponses, SetGatewayResponse(gateway))
	}
	return ctx.JSON(http.StatusOK, gatewayResponses)
}

func (h *Handler) AdminGetGateway(ctx echo.Context) error {
	//TODO without authorize
	return nil
}

func (h *Handler) AdminUpdateGateway(ctx echo.Context) error {
	//TODO block unblock,...
	return nil
}

// transaction handlers

func (h *Handler) AdminListTransactions(ctx echo.Context) error {
	transactions, err := services.ListAllTransaction(h.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server error in getting transactions")
	}
	var transactionResponses []TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, SetVerifyTransactionResponse(transaction))
	}
	return ctx.JSON(http.StatusOK, transactionResponses)
}

func (h *Handler) AdminGetTransaction(ctx echo.Context) error {
	//TODO
	return nil
}

func (h *Handler) AdminUpdateTransaction(ctx echo.Context) error {
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
