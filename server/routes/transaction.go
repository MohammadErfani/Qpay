package routes

import (
	"Qpay/database"
	"Qpay/server/handlers"
	"github.com/labstack/echo/v4"
)

func TransactionGroup(ctx *echo.Group) {
	tr := &handlers.TransactionHandler{
		UserID: 1,
		DB:     database.DB(),
	}
	ctx.GET("/transaction/list", tr.ListAllTransaction)
	ctx.GET("/transaction/find/:id", tr.FindTransaction)
	ctx.POST("/transaction/filter", tr.FilterTransaction)
	ctx.POST("/transaction/search", tr.SearchTransaction)
	// یکی کردن پرسنال و بیزینس
	ctx.POST("/transaction/create/:route", tr.CreateTransaction)
	//ctx.POST("/gateway/:route", tr.RequestPersonalTransaction)
	//ctx.POST("/transaction/PaymentRequest", tr.RequestBusinessTransaction)
	// get : دریافت اطلاعات تراکنش
	ctx.GET("/transaction/StartPay/:id", tr.GetTransactionForStart)
	ctx.POST("/transaction/StartPay", tr.BeginTransaction)
	ctx.POST("/transaction/PaymentVerification", tr.VerifyTransaction)
}
