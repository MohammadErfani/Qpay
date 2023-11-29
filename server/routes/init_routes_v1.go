package routes

import (
	"Qpay/config"
	_ "Qpay/docs"
	"Qpay/server/middlewares"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"gorm.io/gorm"
)

// InitRoutesV1 version 1 api routes:
func InitRoutesV1(db *gorm.DB, cfg *config.Config) *echo.Echo {

	authMiddleware := &middlewares.Auth{
		JWT: &cfg.JWT,
	}

	e := echo.New()
	e.GET("/doc/*", echoSwagger.WrapHandler)
	v1 := e.Group("/api/v1")
	v1Auth := v1.Group("", authMiddleware.AuthMiddleware)
	v1Admin := v1.Group("", authMiddleware.AuthMiddleware, middlewares.AdminMiddleware)
	UserGroup(v1)
	BankAccountGroup(v1Auth)
	AdminGroup(v1Admin)
	GatewayGroup(v1Auth)
	AuthGroup(v1, db, cfg)
	TransactionGroup(v1Auth)
	PaymentGroup(v1)
	return e
}
