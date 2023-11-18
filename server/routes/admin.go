package routes

import (
	"Qpay/database"
	"Qpay/server/handlers"
	"github.com/labstack/echo/v4"
)

func AdminGroup(adminG *echo.Group) {
	aH := handlers.AdminHandler{
		DB: database.DB(),
	}
	adminG.POST("/admin/register", aH.AdminCreate)
}
