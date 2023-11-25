package handlers

import (
	"Qpay/models"
	"Qpay/services"
	"Qpay/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type TransactionHandler struct {
	DB     *gorm.DB
	UserID uint
}

func (tr *TransactionHandler) ListAllTransaction(ctx echo.Context) error {
	return nil
}

func (tr *TransactionHandler) FindTransaction(ctx echo.Context) error {
	return nil
}

func (tr *TransactionHandler) FilterTransaction(ctx echo.Context) error {
	//	امکان فیلتر کردن تراکنش‌ها بر حسب تاریخ و یا قیمت (بازه زمانی و یا قیمتی)
	return nil
}
func (tr *TransactionHandler) SearchTransaction(ctx echo.Context) error {
	//	امکان جستجو در تراکنش‌های ثبت شده بر حسب تاریخ و یا قیمت (بازه زمانی و یا قیمتی)
	return nil
}

func (tr *TransactionHandler) CreateTransaction(ctx echo.Context) error {
	// دریافت پست مقادیر زیر
	// گرفتن درگاه بر اساس route
	//	مقدار پرداخت
	// ست کردن مبلغ commission
	//	شماره موبایل

	//	ریسپانس مقادیر زیر
	//	آی دی تراکنش
	return nil
}

//	func (tr *TransactionHandler) RequestPersonalTransaction(ctx echo.Context) error {
//		// دریافت پست مقادیر زیر
//		//	آدرس درگاه - Route - باید به آی دی تبدیل بشه
//		//	مقدار پرداخت
//		//	شماره موبایل
//
//		//	ریسپانس مقادیر زیر
//		//	آی دی تراکنش
//		return nil
//	}
//
//	func (tr *TransactionHandler) RequestBusinessTransaction(ctx echo.Context) error {
//		// دریافت پست مقادیر زیر
//		//	آی دی درگاه merchantId
//		//	مقدار پرداخت
//		//	شماره موبایل
//
//		//	ریسپانس مقادیر زیر
//		//	آی دی تراکنش
//		return nil
//	}
func (tr *TransactionHandler) GetTransactionForStart(ctx echo.Context) error {
	// response:

	return nil
}

type PaymentTransactionRequest struct {
	PaymentAmount       float64 `json:"payment_amount"`
	PurchaserCard       string  `json:"purchaser_card"`
	CardMonth           int     `json:"card_month"`
	CardYear            int     `json:"card_year"`
	PhoneNumber         string  `json:"phone_number"`
	TransactionID       uint    `json:"transaction_id"`
	PaymentConfirmation bool    `json:"payment_confirmation"` //	دستور پرداخت و کم کردن موجودی (کنسل تراکنش - پرداخت)
}
type PaymentTransactionResponse struct {
	TransactionID  uint    `json:"transaction_id"`
	TrackingCode   string  `json:"tracking_code"`
	Status         uint8   `json:"status"`
	PaymentAmouont float64 `json:"payment_amount"`
	PurchaserCard  string  `json:"purchaser_card"`
}

func BeginTransactionResponse(transaction models.Transaction) PaymentTransactionResponse {
	return PaymentTransactionResponse{
		TransactionID:  transaction.ID,
		TrackingCode:   transaction.TrackingCode,
		Status:         transaction.Status,
		PaymentAmouont: transaction.PaymentAmount,
		PurchaserCard:  transaction.PurchaserCard,
	}
}
func (tr *TransactionHandler) BeginTransaction(ctx echo.Context) error {
	var req PaymentTransactionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	// اینجا چک میکنه اگه طرف فیلد پیمنت کانفیرم رو فالس داده بود
	//یعنی میخواد پرداخت رو کنسل کنه و پرداخت انجام نده
	if !req.PaymentConfirmation {
		if err := services.CancelledTransaction(tr.DB, req.TransactionID); err != nil {
			return ctx.JSON(http.StatusBadRequest, err.Error())
		}
		return ctx.JSON(http.StatusNotAcceptable, "your Payment Transaction is Canceled")
	}

	// حالا که همه چیز آماده انجام تراکنش هست باید اول بررسی شود که
	// فلیدهای لازم درون درخواست وجود داند یا خیر؟
	if err := ValidateTransaction(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	// اینجا باید به ماک متصل بشم و یه خروجی ازش بگیرم که مثلا از کارت مشتری پول کم شده
	transaction, err := services.CreateTransaction(tr.DB, req.TransactionID, req.PaymentAmount, req.CardYear, req.CardMonth, req.PhoneNumber, req.PurchaserCard)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, BeginTransactionResponse(*transaction))

}
func ValidateTransaction(gateway *PaymentTransactionRequest) error {

	requiredFields := map[string]string{
		"purchaser_card": gateway.PurchaserCard,
		"phone_number":   gateway.PhoneNumber,
	}
	requiredFieldsInt := map[string]int{
		"card_month": gateway.CardMonth,
		"card_year":  gateway.CardYear,
	}
	requiredFieldsUint := map[string]uint{
		"tracking_code": gateway.TransactionID,
	}
	requiredFieldsFloat64 := map[string]float64{
		"payment_amount": gateway.PaymentAmount,
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

	if err := utils.IsRequiredFloat64(requiredFieldsFloat64); err != nil {
		return err
	}
	return nil
}

type TransactionResponse struct {
	TrackingCode  string  `json:"tracking_code"`
	Status        string  `json:"status"`
	PurchaserCard string  `json:"purchaser_card"`
	PaymentAmount float64 `json:"payment_amount"`
	PhoneNumber   string  `json:"phone_number"`
	PaymentDate   string  `json:"payment_date"`
}

func (tr *TransactionHandler) VerifyTransaction(ctx echo.Context) error {
	var transaction models.Transaction
	trackingCode := ctx.Param("tracking_code")
	transaction, err := services.GetSpecificTransaction(tr.DB, trackingCode)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Transaction does not exist!")
	}
	return ctx.JSON(http.StatusOK, SetVerifyTransactionResponse(transaction))
}
func SetVerifyTransactionResponse(transaction models.Transaction) TransactionResponse {
	var status string
	if transaction.Status == models.NotPaid {
		status = "NotPaid"
	} else if transaction.Status == models.NotSuccessfully {
		status = "NotSuccessfully"
	} else if transaction.Status == models.IssueOccurred {
		status = "IssueOccurred"
	} else if transaction.Status == models.Blocked {
		status = "Blocked"
	} else if transaction.Status == models.Refund {
		status = "Refund"
	} else if transaction.Status == models.Cancelled {
		status = "Cancelled"
	} else if transaction.Status == models.ReturnToGateway {
		status = "ReturnToGateway"
	} else if transaction.Status == models.AwaitingConfirmation {
		status = "AwaitingConfirmation"
	} else if transaction.Status == models.Confirmed {
		status = "Confirmed"
	} else if transaction.Status == models.Paid {
		status = "Paid"
	}

	return TransactionResponse{
		TrackingCode:  transaction.TrackingCode,
		Status:        status,
		PurchaserCard: transaction.PurchaserCard,
		PaymentAmount: transaction.PaymentAmount,
		PhoneNumber:   transaction.PhoneNumber,
		PaymentDate:   transaction.CreatedAt.Format("2006-01-02 15:04"),
	}
}
