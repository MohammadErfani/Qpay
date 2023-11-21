package routes

import (
	"Qpay/database"
	"Qpay/server/handlers"
	"github.com/labstack/echo/v4"
)

func TransactionGroup(ctx echo.Group) {
	tr := &handlers.TransHandler{
		UserID: 1,
		DB:     database.DB(),
	}
	_ = tr
}
