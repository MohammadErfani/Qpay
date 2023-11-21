package handlers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TransactionRequest struct {
	GatewayID     uint    `json:"gateway_id"`
	PaymentAmount float64 `json:"payment_amount"`
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

func (tr *TransactionHandler) RequestPersonalTransaction(ctx echo.Context) error {
	// دریافت پست مقادیر زیر
	//	آدرس درگاه - Route
	//	مقدار پرداخت
	//	شماره موبایل

	//	ریسپانس مقادیر زیر
	//	آی دی تراکنش
	return nil
}

func (tr *TransactionHandler) RequestBusinessTransaction(ctx echo.Context) error {
	// دریافت پست مقادیر زیر
	//	آی دی درگاه
	//	مقدار پرداخت
	//	شماره موبایل

	//	ریسپانس مقادیر زیر
	//	آی دی تراکنش
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
	return nil
}
