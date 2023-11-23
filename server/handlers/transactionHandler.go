package handlers

import (
	"Qpay/models"
	"Qpay/services"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type TransactionResponse struct {
	TrackingCode  string  `json:"tracking_code"`
	Status        string  `json:"status"`
	PurchaserCard string  `json:"purchaser_card"`
	PaymentAmount float64 `json:"payment_amount"`
	PhoneNumber   string  `json:"phone_number"`
}
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
func (tr *TransactionHandler) BeginTransaction(ctx echo.Context) error {
	// دریافت پست مقادیر زیر
	// شماره تراکنش
	//	شماره کارت
	//	شماره cvv2
	//	رمز کارت بانکی
	//	ماه انقضا کارت
	//	سال انقضا کارت
	//	استاتوس (کنسل تراکنش - تایید پرداخت)
	//	آدرس بازگشتی  callbackUrl

	//	ریسپانس مقادیر زیر
	//	آی دی تراکنش
	//	استاتوس تراکنش
	//	مبلغ پرداخت شده
	//	کپی و پیست آدرس بازگشتی
	//	۴ رقم آخر شماره کارت - یا برای ساده تر شدن کل شماره کارت
	return nil

}

func (tr *TransactionHandler) VerifyTransaction(ctx echo.Context) error {
	// دریافت مقادیر زیر جهت تایید وضعیت تراکنش
	//	شماره تراکنش

	//	ریسپانس مقادیر بازگشتی
	//	وضعیت تراکنش
	//	۴ رقم آخر شماره کارت - یا برای سایده تر شدن کل شماره کارت
	//	تاریخ و ساعت کم شدن پول
	//	مبلغ پرداخت شده

	var transaction models.Transaction
	trackingCode := ctx.Param("tracking_code")
	transaction, err := services.GetSpecificTransaction(tr.DB, trackingCode)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Transaction does not exist!")
	}
	return ctx.JSON(http.StatusOK, SetTransactionResponse(transaction))
}
func SetTransactionResponse(transaction models.Transaction) TransactionResponse {
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
	}
}
