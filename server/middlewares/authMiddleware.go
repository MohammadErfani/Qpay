package middlewares

import (
	"Qpay/config"
	"Qpay/utils"
	"fmt"
	"net/http"
	"strings"

	echo "github.com/labstack/echo/v4"
)

const (
	AuthHeader         = "Authorization"
	Bearer             = "bearer"
	UserIdContextField = "user_id"
)

type Auth struct {
	JWT *config.JWT
}

func (a *Auth) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		tokenHeaeder := ctx.Request().Header.Get("Authorization")
		if tokenHeaeder == "" {
			return ctx.NoContent(http.StatusUnauthorized)
		}

		tokenParam := strings.Split(tokenHeaeder, " ")
		if len(tokenParam) < 2 {
			return ctx.NoContent(http.StatusUnauthorized)
		}

		tokenType := strings.ToLower(tokenParam[0])
		if tokenType != Bearer {
			return ctx.NoContent(http.StatusUnauthorized)
		}

		token := tokenParam[1]

		credential, err := utils.VerifyToken(a.JWT, token)

		if err != nil {

			fmt.Printf("ERROR IN VERIFY:: %v\n", err)
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"hasError": "true", "message": "token not valid"})
		}

		ctx.Set(UserIdContextField, credential.ID)
		return next(ctx)
	}
}
