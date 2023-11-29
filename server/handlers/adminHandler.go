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
	PercentPerTrans float64 `json:"percent_per_transaction"`
}

type StatusRequest struct {
	Status string `json:"status"`
}

type CommissionResponse struct {
	ID              uint    `json:"id"`
	AmountPerTrans  float64 `json:"amount_per_transaction"`
	PercentPerTrans float64 `json:"Percent_per_transaction"`
	Status          string  `json:"status"`
}

// AdminCreate godoc
// @Summary Create a new admin
// @Description Register a new admin with the provided details.
// @Tags admin
// @Accept json
// @Produce json
// @Param adminRequest body AdminRequest true "Admin details"
// @Success 201 {object} Response "Success message"
// @Failure 400 {object} Response "{"status": "error", "message": "Bind Error"}"
// @Failure 403 {object} Response "{"status": "error", "message": "Invalid admin details"}"
// @Failure 409 {object} Response "{"status": "error", "message": "Username or email already exists"}"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error in create admin"}"
// @Security ApiKeyAuth
// @Router /admin/register [post]
func (h *Handler) AdminCreate(ctx echo.Context) error {
	var req AdminRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Bind Error",
		})
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
		return ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Internal server error in create admin",
		})
	}
	return ctx.JSON(http.StatusCreated, Response{Status: "success", Message: "Admin created successfully"})

}

// AdminListAllCommission godoc
// @Summary List all commissions
// @Description Retrieve a list of all commissions.
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {array} CommissionResponse "List of commissions"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error in getting commissions"}"
// @Security ApiKeyAuth
// @Router /admin/commission [get]
func (h *Handler) AdminListAllCommission(ctx echo.Context) error {
	commissions, err := services.ListAllCommission(h.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Internal server error in getting commissions",
		})
	}
	var commResponses []CommissionResponse
	for _, comm := range commissions {
		commResponses = append(commResponses, SetCommissionsResponse(comm))
	}
	return ctx.JSON(http.StatusOK, commResponses)

}

// AdminCreateCommission godoc
// @Summary Create a new commission
// @Description Create a new commission with the provided details.
// @Tags admin
// @Accept json
// @Produce json
// @Param commissionRequest body CommissionRequest true "Commission details"
// @Success 201 {object} CommissionResponse "Created commission details"
// @Failure 400 {object} Response "{"status": "error", "message": "Bind Error"}"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error in creating commission"}"
// @Security ApiKeyAuth
// @Router /admin/commission [post]
func (h *Handler) AdminCreateCommission(ctx echo.Context) error {
	var req CommissionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Bind Error",
		})
	}
	comm, err := services.CreateCommission(h.DB, req.AmountPerTrans, req.PercentPerTrans)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Internal server error in creating commission",
		})
	}
	return ctx.JSON(http.StatusCreated, SetCommissionsResponse(*comm))
}

// AdminGetCommission godoc
// @Summary Get commission by ID
// @Description Retrieve details of a specific commission by its ID.
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "Commission ID"
// @Success 200 {object} CommissionResponse "Commission details"
// @Failure 404 {object} Response "{"status": "error", "message": "Commission not found!"}"
// @Security ApiKeyAuth
// @Router /admin/commission/{id} [get]
func (h *Handler) AdminGetCommission(ctx echo.Context) error {
	commID := ctx.Param("id")
	comm, err := services.GetCommission(h.DB, "id", fmt.Sprintf("%v", commID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "commission not found!",
		})
	}
	return ctx.JSON(http.StatusOK, SetCommissionsResponse(*comm))
}

// AdminListUsers godoc
// @Summary List all users
// @Description Retrieve a list of all users.
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {array} UserResponse "List of users"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error in getting users"}"
// @Security ApiKeyAuth
// @Router /admin/user [get]
func (h *Handler) AdminListUsers(ctx echo.Context) error {
	users, err := services.ListAllUser(h.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Internal server error in getting users",
		})
	}
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, SetUsersResponse(user))
	}
	return ctx.JSON(http.StatusOK, userResponses)
}

// AdminGetUser godoc
// @Summary Get user by ID
// @Description Retrieve details of a specific user by its ID.
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse "User details"
// @Failure 404 {object} Response "{"status": "error", "message": "User not found!"}"
// @Security ApiKeyAuth
// @Router /admin/user/{id} [get]
func (h *Handler) AdminGetUser(ctx echo.Context) error {
	userID := ctx.Param("id")
	user, err := services.GetUser(h.DB, "id", fmt.Sprintf("%v", userID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "user not found!",
		})
	}
	return ctx.JSON(http.StatusOK, SetUsersResponse(*user))
}

// AdminUpdateUser godoc
// @Summary Update user status by ID
// @Description Update the status of a specific user by its ID.
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param updateRequest body StatusRequest true "Updated status details"
// @Success 200 {object} Response "Success message"
// @Failure 400 {object} Response "{"status": "error", "message": "Bind error"}"
// @Failure 400 {object} Response "{"status": "error", "message": "Status field is unsupported"}"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error"}"
// @Security ApiKeyAuth
// @Router /admin/user/{id} [patch]
func (h *Handler) AdminUpdateUser(ctx echo.Context) error {
	userID := ctx.Param("id")
	user, err := services.GetUser(h.DB, "id", fmt.Sprintf("%v", userID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "user not found!",
		})
	}
	var req StatusRequest
	if err = ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Bind error",
		})
	}
	err = services.ChangeUserGatewaysStatus(h.DB, user, req.Status)
	if err != nil {
		if err.Error() == "status field is unsupported" {
			return ctx.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "internal server error",
		})
	}
	return ctx.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "user gateways updated successfully",
	})

}

// AdminListAllGateways godoc
// @Summary List all gateways
// @Description Retrieve a list of all gateways.
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {array} GatewayResponse "List of gateways"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error in getting gateways"}"
// @Security ApiKeyAuth
// @Router /admin/gateway [get]
func (h *Handler) AdminListAllGateways(ctx echo.Context) error {
	gateways, err := services.ListAllGateway(h.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Internal server error in getting gateways",
		})
	}
	var gatewayResponses []GatewayResponse
	for _, gateway := range gateways {
		gatewayResponses = append(gatewayResponses, SetGatewayResponse(gateway))
	}
	return ctx.JSON(http.StatusOK, gatewayResponses)
}

// AdminGetGateway godoc
// @Summary Get gateway by ID
// @Description Retrieve details of a specific gateway by its ID.
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "Gateway ID"
// @Success 200 {object} GatewayResponse "Gateway details"
// @Failure 404 {object} Response "{"status": "error", "message": "Gateway not found!"}"
// @Security ApiKeyAuth
// @Router /admin/gateway/{id} [get]
func (h *Handler) AdminGetGateway(ctx echo.Context) error {
	gatewayID := ctx.Param("id")
	gateway, err := services.GetGateway(h.DB, "id", fmt.Sprintf("%v", gatewayID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "gateway not found!",
		})
	}
	return ctx.JSON(http.StatusOK, SetGatewayResponse(*gateway))
}

// AdminUpdateGateway godoc
// @Summary Update gateway status by ID
// @Description Update the status of a specific gateway by its ID.
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "Gateway ID"
// @Param updateRequest body StatusRequest true "Updated status details"
// @Success 200 {object} GatewayResponse "Updated gateway details"
// @Failure 400 {object} Response "{"status": "error", "message": "Bind error"}"
// @Failure 400 {object} Response "{"status": "error", "message": "Status field is unsupported"}"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error"}"
// @Security ApiKeyAuth
// @Router /admin/gateway/{id} [patch]
func (h *Handler) AdminUpdateGateway(ctx echo.Context) error {
	gatewayID := ctx.Param("id")
	gateway, err := services.GetGateway(h.DB, "id", fmt.Sprintf("%v", gatewayID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "gateway not found!",
		})
	}
	var req StatusRequest
	if err = ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Bind error",
		})
	}
	gateway, err = services.SetStatusGateway(h.DB, gateway, req.Status)
	if err != nil {
		if err.Error() == "status field is unsupported" {
			return ctx.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "internal server error",
		})
	}
	return ctx.JSON(http.StatusOK, SetGatewayResponse(*gateway))
}

// transaction handlers

// AdminListTransactions godoc
// @Summary List all transactions
// @Description Retrieve a list of all transactions.
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {array} TransactionResponse "List of transactions"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error in getting transactions"}"
// @Security ApiKeyAuth
// @Router /admin/transaction [get]
func (h *Handler) AdminListTransactions(ctx echo.Context) error {
	transactions, err := services.ListAllTransaction(h.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Internal server error in getting transactions",
		})
	}
	var transactionResponses []TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, SetTransactionResponse(transaction))
	}
	return ctx.JSON(http.StatusOK, transactionResponses)
}

// AdminGetTransaction godoc
// @Summary Get transaction by ID
// @Description Retrieve details of a specific transaction by its ID.
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} TransactionResponse "Transaction details"
// @Failure 404 {object} Response "{"status": "error", "message": "Transaction not found!"}"
// @Security ApiKeyAuth
// @Router /admin/transaction/{id} [get]
func (h *Handler) AdminGetTransaction(ctx echo.Context) error {
	transactionID := ctx.Param("id")
	transaction, err := services.GetTransaction(h.DB, "id", fmt.Sprintf("%v", transactionID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "transaction not found!",
		})
	}
	return ctx.JSON(http.StatusOK, SetTransactionResponse(*transaction))
}

// AdminUpdateTransaction godoc
// @Summary Update transaction status by ID
// @Description Update the status of a specific transaction by its ID.
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param updateRequest body StatusRequest true "Updated status details"
// @Success 200 {object} TransactionResponse "Updated transaction details"
// @Failure 400 {object} Response "{"status": "error", "message": "Bind error"}"
// @Failure 400 {object} Response "{"status": "error", "message": "Status field is unsupported"}"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error"}"
// @Security ApiKeyAuth
// @Router /admin/transaction/{id} [patch]
func (h *Handler) AdminUpdateTransaction(ctx echo.Context) error {
	transactionID := ctx.Param("id")
	transaction, err := services.GetTransaction(h.DB, "id", fmt.Sprintf("%v", transactionID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "transaction not found!",
		})
	}
	var req StatusRequest
	if err = ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Bind error",
		})
	}
	transaction, err = services.SetStatusTransaction(h.DB, transaction, req.Status)
	if err != nil {
		if err.Error() == "status field is unsupported" {
			return ctx.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "internal server error",
		})
	}
	return ctx.JSON(http.StatusOK, SetTransactionResponse(*transaction))
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
