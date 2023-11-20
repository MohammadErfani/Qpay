package routes

import (
	"Qpay/config"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// InitRoutesV1 version 1 api routes:
func InitRoutesV1(db *gorm.DB, cfg *config.Config) *echo.Echo {

	e := echo.New()

	v1 := e.Group("/api/v1")
	v1.GET("/test", func(ctx echo.Context) error {
		return nil
	})

	UserGroup(v1)
	AuthGroup(v1, db, cfg)
	return e
}
