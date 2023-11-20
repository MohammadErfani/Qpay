package routes

import (
	// "Qpay/server/middlewares"
	// "errors"
	// "net/http"

	// "github.com/golang-jwt/jwt/v5"
	// "github.com/labstack/echo-jwt/v4"

	"Qpay/config"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// InitRoutesV1 version 1 api routes:
func InitRoutesV1(db *gorm.DB, cfg *config.Config) *echo.Echo {

	e := echo.New()

	v1 := e.Group("/api/v1")

	// e.Use(echojwt.WithConfig(echojwt.Config{
	// 	SigningKey: []byte("secret"),
	// }))

	// test
	v1.GET("/test", func(ctx echo.Context) error {
		// v1.POST("/login")
		// token, ok := ctx.Get("user").(*jwt.Token) // key for example at this context is user
		// if !ok {
		// 	return errors.New("JWT token missing or invalid")
		// }
		// claims, ok := token.Claims.(jwt.MapClaims)

		// if !ok {
		// 	return errors.New("failed to cast claims as jwt.MapClaims")
		// }

		// return ctx.JSON(http.StatusOK, claims)
		return nil
	})
	UserGroup(v1)
	AuthGroup(v1, db, cfg)
	return e
}
