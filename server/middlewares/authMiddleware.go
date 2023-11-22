package middlewares

import (
  "net/http"
	"Qpay/utils"

	echo "github.com/labstack/echo/v4"
)

// type Auth struct {
// 	DB  *gorm.DB
// 	JWT *config.JWT
// }

func AuthMiddleware(c echo.HandlerFunc) echo.HandlerFunc{
  return func(c echo.Context) error {
    token := c.Request().Header.Get("Authorization")
    checkToken := utils.ValidationToken(token)
    if checkToken != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"hasError": "true", "message": "Your token not valid or expired!"})
    }
return nil
  }
}
