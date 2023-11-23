package routes

import (
	"Qpay/config"
	"Qpay/server/handlers"
	"Qpay/server/middlewares"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// InitRoutesV1 version 1 api routes:
func InitRoutesV1(db *gorm.DB, cfg *config.Config) *echo.Echo {

	authMiddleware := &middlewares.Auth{
		JWT: &cfg.JWT,
	}

	e := echo.New()

	v1 := e.Group("/api/v1")
	v1.GET("/test", handlers.TestHandler, authMiddleware.AuthMiddleware)

	UserGroup(v1)
	AuthGroup(v1, db, cfg)
	return e
}
