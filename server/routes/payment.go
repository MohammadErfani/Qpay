package routes

import (
	"Qpay/database"
	"Qpay/server/handlers"
	"github.com/labstack/echo/v4"
)

func PaymentGroup(ctx *echo.Group) {
	h := handlers.Handler{
		DB: database.DB(),
	}
	// یکی کردن پرسنال و بیزینس
	ctx.POST("/transaction/create/:route", h.CreateTransaction)
	// get : دریافت اطلاعات تراکنش
	ctx.GET("/transaction/StartPay/:id", h.GetTransactionForStart)
	ctx.POST("/transaction/StartPay", h.BeginTransaction)
	ctx.POST("/transaction/PaymentVerification", h.VerifyTransaction)
}
