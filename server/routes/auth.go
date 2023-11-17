package routes

import (
	"Qpay/utils"
	"encoding/json"
	// "fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

type User struct {
	email       string `json:"email"`
	phoneNumber string `json:"phoneNumber"`
	password    string `json:"password"`
}

func AuthGroup(authG *echo.Group) {
	json_map := make(map[string]interface{})

	// user := new(User)
	authG.POST("/login", func(ctx echo.Context) error {

		err := json.NewDecoder(ctx.Request().Body).Decode(&json_map)

		email := json_map["email"].(string)
		emailError := utils.IsValidEmail(email)
		if emailError != nil {
			return ctx.JSON(http.StatusBadRequest, emailError)
		}

		phoneNumber := json_map["phoneNumber"].(string)
		phoneNumberError := utils.IsValidPhoneNumber(phoneNumber)
		if phoneNumberError != nil {
			return ctx.JSON(http.StatusBadRequest, phoneNumberError)
		}
		// password := json_map["password"]

		if err != nil {
			return ctx.String(http.StatusBadRequest, "Bad request")
		}

		// user.email = email.String()

		return ctx.JSON(http.StatusOK, json_map)
	})
}
