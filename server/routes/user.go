package routes

import (
	"Qpay/server/handlers"
	echo "github.com/labstack/echo/v4"
)

func UserGroup(userG *echo.Group) {
	userG.POST("/register", handlers.CreateUser)
}
