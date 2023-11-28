package routes

import (
	"Qpay/database"
	"Qpay/server/handlers"
	echo "github.com/labstack/echo/v4"
)

func UserGroup(userG *echo.Group) {
	h := &handlers.Handler{
		DB: database.DB(),
	}
	userG.POST("/register", h.CreateUser)
}
