package routes

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

// InitRoutesV1 version 1 api routes:
func InitRoutesV1() *echo.Echo {
	e := echo.New()

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secret"),
	}))

	v1 := e.Group("/api/v1")
	// test
	v1.GET("/test", func(ctx echo.Context) error {
		token, ok := ctx.Get("user").(*jwt.Token) // key for example at this context is user
		if !ok {
			return errors.New("JWT token missing or invalid")
		}
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			return errors.New("failed to cast claims as jwt.MapClaims")
		}

		return ctx.JSON(http.StatusOK, claims)
	})
	UserGroup(v1)
	AuthGroup(v1)
	return e
}
