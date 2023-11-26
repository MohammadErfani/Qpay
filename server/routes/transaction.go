package routes

import (
	"Qpay/database"
	"Qpay/server/handlers"
	"github.com/labstack/echo/v4"
)

func TransactionGroup(ctx *echo.Group) {
	h := &handlers.Handler{
		UserID: 1,
		DB:     database.DB(),
	}
	ctx.GET("/transaction/list/", h.ListAllTransaction)
	ctx.GET("/transaction/find/:id", h.FindTransaction)
	ctx.POST("/transaction/filter/", h.FilterTransaction)
	ctx.POST("/transaction/search/", h.SearchTransaction)

}
