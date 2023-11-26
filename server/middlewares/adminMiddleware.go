package middlewares

import (
	"Qpay/database"
	"Qpay/services"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := ctx.Get(UserIdContextField).(int)
		fmt.Println(userID)
		if err := services.CheckIsAdmin(database.DB(), uint(userID)); err != nil {
			if err.Error() == "unAuthorize" {
				return ctx.NoContent(http.StatusUnauthorized)
			}
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return next(ctx)
	}
}
