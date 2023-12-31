package handlers

import (
	"Qpay/models"
	"Qpay/services"
	"Qpay/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CreateTransactionRequest struct {
	PhoneNumber   string  `json:"phone_number"`
	PaymentAmount float64 `json:"payment_amount"`
}
type CreateTransactionResponse struct {
	TransactionID uint `json:"id"`
}

type TransactionStartResponse struct {
	PaymentAmount float64 `json:"payment_amount"`
	OwnerName     string  `json:"owner_name"`
}
type TransactionResponse struct {
	TrackingCode  string  `json:"tracking_code"`
	Status        string  `json:"status"`
	PurchaserCard string  `json:"purchaser_card"`
	PaymentAmount float64 `json:"payment_amount"`
	PhoneNumber   string  `json:"phone_number"`
	PaymentDate   string  `json:"payment_date"`
}

type FilterRequest struct {
	StartDate *string  `json:"start_date"`
	EndDate   *string  `json:"end_date"`
	MinAmount *float64 `json:"min_amount"`
	MaxAmount *float64 `json:"max_amount"`
}

// ListAllTransaction godoc
// @Summary List all transactions for the authenticated user
// @Description Retrieve a list of all transactions associated with the authenticated user.
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {array} TransactionResponse "List of transactions"
// @Failure 400 {object} Response "{"status": "error", "message": "You don't have any transaction"}"
// @Security ApiKeyAuth
// @Router /transaction/list [get]
func (h *Handler) ListAllTransaction(ctx echo.Context) error {
	h.SetUserID(ctx)
	transactions, err := services.GetUserTransactions(h.DB, h.UserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "You don't have Any Transaction",
		})
	}
	var TransactionResponses []TransactionResponse
	for _, ba := range transactions {
		TransactionResponses = append(TransactionResponses, SetTransactionResponse(ba))
	}
	return ctx.JSON(http.StatusOK, TransactionResponses)
}

// FindTransaction godoc
// @Summary Find a transaction by ID for the authenticated user
// @Description Retrieve details of a specific transaction associated with the authenticated user.
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} TransactionResponse "Transaction details"
// @Failure 400 {object} Response "{"status": "error", "message": "Transaction ID is not correct"}"
// @Failure 404 {object} Response "{"status": "error", "message": "Transaction does not exist!"}"
// @Security ApiKeyAuth
// @Router /transaction/find/{id} [get]
func (h *Handler) FindTransaction(ctx echo.Context) error {
	h.SetUserID(ctx)
	var transaction models.Transaction
	transactionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "gateway is not correct",
		})
	}
	transaction, err = services.FindTransaction(h.DB, h.UserID, uint(transactionID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Gateway does not exist!",
		})
	}
	return ctx.JSON(http.StatusOK, SetTransactionResponse(transaction))
}

// FilterTransaction godoc
// @Summary Filter transactions for the authenticated user
// @Description Filter transactions based on start date, end date, and payment amount range.
// @Tags transactions
// @Accept json
// @Produce json
// @Param filterRequest body FilterRequest true "Filter criteria"
// @Success 200 {array} TransactionResponse "List of filtered transactions"
// @Failure 400 {object} Response "{"status": "error", "message": "Bind Error"}"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error"}"
// @Security ApiKeyAuth
// @Router /transaction/filter [post]
func (h *Handler) FilterTransaction(ctx echo.Context) error {
	h.SetUserID(ctx)
	var filter FilterRequest
	if err := ctx.Bind(&filter); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Bind Error",
		})
	}
	filtered, err := services.FilterTransaction(h.DB, h.UserID, filter.StartDate, filter.EndDate, filter.MinAmount, filter.MaxAmount)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: err.Error(),
		})
	}
	var filteredResponse []TransactionResponse
	for _, i := range filtered {
		filteredResponse = append(filteredResponse, TransactionResponse{
			TrackingCode:  i.TrackingCode,
			Status:        GetStatus(uint(i.Status)),
			PurchaserCard: i.PurchaserCard,
			PaymentAmount: i.PaymentAmount,
			PhoneNumber:   i.PhoneNumber,
			PaymentDate:   i.CreatedAt.Format("2006-01-02 15:04"),
		})
	}
	return ctx.JSON(http.StatusOK, filteredResponse)
}

// CreateTransaction godoc
// @Summary Create a new transaction for the authenticated user
// @Description Create a new transaction with the provided payment details.
// @Tags transactions
// @Accept json
// @Produce json
// @Param route path string true "Gateway route"
// @Param transactionRequest body CreateTransactionRequest true "Transaction details"
// @Success 200 {object} CreateTransactionResponse "Transaction ID"
// @Failure 400 {object} Response "{"status": "error", "message": "Bind Error"}"
// @Failure 404 {object} Response "{"status": "error", "message": "No gateway with such route"}"
// @Failure 403 {object} Response "{"status": "error", "message": "SHEBA doesn't match your credentials"}"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error in create transaction"}"
// @Security ApiKeyAuth
// @Router /transaction/create/{route} [post]
func (h *Handler) CreateTransaction(ctx echo.Context) error {
	route := ctx.Param("route")
	var req CreateTransactionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Bind Error",
		})
	}
	gateway, err := services.GetGateway(h.DB, "route", route)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "No get way with such route",
		})
	}

	commission := req.PaymentAmount*gateway.Commission.PercentagePerTrans + gateway.Commission.AmountPerTrans
	model, err := services.CreateTransaction(h.DB, gateway.ID, req.PaymentAmount, req.PhoneNumber, commission)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, CreateTransactionResponse{TransactionID: model.ID})
}

// GetTransactionForStart godoc
// @Summary Get transaction details for starting a transaction
// @Description Get transaction details to start a payment transaction.
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} TransactionStartResponse "Transaction details for start"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error"}"
// @Security ApiKeyAuth
// @Router /transaction/start/{id} [get]
func (h *Handler) GetTransactionForStart(ctx echo.Context) error {
	transactionid, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	transaction, err := services.GetTransactionByID(h.DB, uint(transactionid))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	gateway, err := services.GetGateway(h.DB, "id", fmt.Sprintf("%v", transaction.GatewayID))
	return ctx.JSON(http.StatusOK, TransactionStartResponse{PaymentAmount: transaction.PaymentAmount, OwnerName: gateway.BankAccount.AccountOwner})
}

type PaymentTransactionRequest struct {
	PurchaserCard       string `json:"purchaser_card"`
	CardMonth           int    `json:"card_month"`
	CardYear            int    `json:"card_year"`
	Password            int    `json:"password"`
	CVV2                int    `json:"cvv2"`
	TransactionID       uint   `json:"transaction_id"`
	PaymentConfirmation bool   `json:"payment_confirmation"` //	دستور پرداخت و کم کردن موجودی (کنسل تراکنش - پرداخت)
}
type PaymentTransactionResponse struct {
	TransactionID uint    `json:"transaction_id"`
	TrackingCode  string  `json:"tracking_code"`
	Status        uint8   `json:"status"`
	PaymentAmount float64 `json:"payment_amount"`
	PurchaserCard string  `json:"purchaser_card"`
}

// BeginTransaction godoc
// @Summary Begin a payment transaction
// @Description Begin a payment transaction with the provided payment details.
// @Tags transactions
// @Accept json
// @Produce json
// @Param transactionRequest body PaymentTransactionRequest true "Payment transaction details"
// @Success 200 {object} PaymentTransactionResponse "Payment transaction details"
// @Failure 400 {object} Response "{"status": "error", "message": "Bind Error"}"
// @Failure 400 {object} Response "{"status": "error", "message": "Invalid transaction details"}"
// @Failure 403 {object} Response "{"status": "error", "message": "Transaction is Canceled"}"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error"}"
// @Security ApiKeyAuth
// @Router /transaction/start [post]
func BeginTransactionResponse(transaction models.Transaction) PaymentTransactionResponse {
	return PaymentTransactionResponse{
		TransactionID: transaction.ID,
		TrackingCode:  transaction.TrackingCode,
		Status:        transaction.Status,
		PaymentAmount: transaction.PaymentAmount,
		PurchaserCard: transaction.PurchaserCard,
	}
}
func (h *Handler) BeginTransaction(ctx echo.Context) error {
	var req PaymentTransactionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	// اینجا چک میکنه اگه طرف فیلد پیمنت کانفیرم رو فالس داده بود
	//یعنی میخواد پرداخت رو کنسل کنه و پرداخت انجام نده
	if !req.PaymentConfirmation {
		if err := services.CancelledTransaction(h.DB, req.TransactionID); err != nil {
			return ctx.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: err.Error(),
			})
		}
		return ctx.JSON(http.StatusNotAcceptable, Response{
			Status:  "error",
			Message: "your Payment Transaction is Canceled",
		})
	}

	// حالا که همه چیز آماده انجام تراکنش هست باید اول بررسی شود که
	// فلیدهای لازم درون درخواست وجود داند یا خیر؟
	if err := ValidateTransaction(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	// اینجا باید به ماک متصل بشم و یه خروجی ازش بگیرم که مثلا از کارت مشتری پول کم شده
	transaction, err := services.PaymentTransaction(h.DB, req.TransactionID, req.CardYear, req.CardMonth, req.PurchaserCard, req.CVV2, req.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, BeginTransactionResponse(transaction))

}
func ValidateTransaction(transaction *PaymentTransactionRequest) error {

	requiredFields := map[string]string{
		"purchaser_card": transaction.PurchaserCard,
	}
	requiredFieldsInt := map[string]int{
		"card_month": transaction.CardMonth,
		"card_year":  transaction.CardYear,
		"password":   transaction.Password,
		"cvv2":       transaction.CVV2,
	}
	requiredFieldsUint := map[string]uint{
		"transaction_id": transaction.TransactionID,
	}
	if err := utils.IsRequired(requiredFields); err != nil {
		return err
	}
	if err := utils.IsRequiredInt(requiredFieldsInt); err != nil {
		return err
	}
	if err := utils.IsRequiredUint(requiredFieldsUint); err != nil {
		return err
	}

	return nil
}

type VerifyTransactionRequest struct {
	TrackingCode string `json:"tracking_code"`
}

// VerifyTransaction godoc
// @Summary Verify a transaction by tracking code
// @Description Verify a transaction by tracking code and retrieve its details.
// @Tags transactions
// @Accept json
// @Produce json
// @Param verifyRequest body VerifyTransactionRequest true "Verification request"
// @Success 200 {object} TransactionResponse "Transaction details"
// @Failure 400 {object} Response "{"status": "error", "message": "Bind Error"}"
// @Failure 404 {object} Response "{"status": "error", "message": "Transaction does not exist!"}"
// @Security ApiKeyAuth
// @Router /transaction/verify [post]
func (h *Handler) VerifyTransaction(ctx echo.Context) error {
	var transaction models.Transaction
	var req VerifyTransactionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "bind error",
		})
	}
	transaction, err := services.GetSpecificTransaction(h.DB, req.TrackingCode)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Transaction does not exist!",
		})
	}
	return ctx.JSON(http.StatusOK, SetTransactionResponse(transaction))
}

func SetTransactionResponse(transaction models.Transaction) TransactionResponse {
	var status string
	status = GetStatus(uint(transaction.Status))
	return TransactionResponse{
		TrackingCode:  transaction.TrackingCode,
		Status:        status,
		PurchaserCard: transaction.PurchaserCard,
		PaymentAmount: transaction.PaymentAmount,
		PhoneNumber:   transaction.PhoneNumber,
		PaymentDate:   transaction.CreatedAt.Format("2006-01-02 15:04"),
	}
}
func GetStatus(statusID uint) string {
	var status string
	if statusID == models.NotPaid {
		status = "NotPaid"
	} else if statusID == models.NotSuccessfully {
		status = "NotSuccessfully"
	} else if statusID == models.IssueOccurred {
		status = "IssueOccurred"
	} else if statusID == models.Blocked {
		status = "Blocked"
	} else if statusID == models.Refund {
		status = "Refund"
	} else if statusID == models.Cancelled {
		status = "Cancelled"
	} else if statusID == models.ReturnToGateway {
		status = "ReturnToGateway"
	} else if statusID == models.AwaitingConfirmation {
		status = "AwaitingConfirmation"
	} else if statusID == models.Confirmed {
		status = "Confirmed"
	} else if statusID == models.Paid {
		status = "Paid"
	}
	return status
}
