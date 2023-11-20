package middlewares

import (
	"Qpay/config"

	echo "github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Auth struct {
	DB  *gorm.DB
	JWT *config.JWT
}

func (a *Auth) AuthMiddleware(next echo.Context) error {
	return nil
}
