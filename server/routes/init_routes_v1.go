package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// InitRoutesV1 version 1 api routes:
func InitRoutesV1() *echo.Echo {
	e := echo.New()
	v1 := e.Group("/api/v1")
	// test
	v1.GET("/test", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "This is Qpay!")
	})
	AuthGroup(v1)
	return e
}
