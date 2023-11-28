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
	ctx.POST("/transaction/create/:route", h.CreateTransaction)
	ctx.GET("/transaction/start/:id", h.GetTransactionForStart)
	ctx.POST("/transaction/start", h.BeginTransaction)
	ctx.POST("/transaction/verify", h.VerifyTransaction)
}
