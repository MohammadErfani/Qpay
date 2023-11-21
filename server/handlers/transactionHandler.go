package handlers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TransactionRequest struct {
	GatewayID     uint    `json:"gateway_id"`
	PaymentAmount float64 `json:"payment_amount"`
}
type TransHandler struct {
	DB     *gorm.DB
	UserID uint
}

func (gh *GatewayHandler) ListAllTransaction(ctx echo.Context) error {
	return nil
}

func (gh *GatewayHandler) FindTransaction(ctx echo.Context) error {
	return nil
}

func (gh *GatewayHandler) FilterTransaction(ctx echo.Context) {
	//	امکان فیلتر کردن تراکنش‌ها بر حسب تاریخ و یا قیمت (بازه زمانی و یا قیمتی)
}
func (gh *GatewayHandler) SearchTransaction(ctx echo.Context) {
	//	امکان جستجو در تراکنش‌های ثبت شده بر حسب تاریخ و یا قیمت (بازه زمانی و یا قیمتی)
}

func (gh *GatewayHandler) RequestPersonalTransaction(ctx echo.Context) {
	// دریافت پست مقادیر زیر
	//	آدرس درگاه - Route
	//	مقدار پرداخت
	//	شماره موبایل

	//	ریسپانس مقادیر زیر
	//	آی دی تراکنش
}

func (gh *GatewayHandler) RequestBusinessTransaction(ctx echo.Context) {
	// دریافت پست مقادیر زیر
	//	آی دی درگاه
	//	مقدار پرداخت
	//	شماره موبایل

	//	ریسپانس مقادیر زیر
	//	آی دی تراکنش
}

func (gh *GatewayHandler) BeginTransaction(ctx echo.Context) {
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

}

func (gh *GatewayHandler) VerifyTransaction(ctx echo.Context) {
	// دریافت مقادیر زیر جهت تایید وضعیت تراکنش
	//	شماره تراکنش

	//	ریسپانس مقادیر بازگشتی
	//	وضعیت تراکنش
	//	۴ رقم آخر شماره کارت - یا برای سایده تر شدن کل شماره کارت
	//	تاریخ و ساعت کم شدن پول
	//	مبلغ پرداخت شده
}
