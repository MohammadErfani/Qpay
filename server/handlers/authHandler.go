package handlers

import (
	"Qpay/config"
	"Qpay/utils"
	"encoding/json"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Auth struct {
	DB  *gorm.DB
	JWT *config.JWT
}

type LoginReq struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginRes struct {
	Token string `json:"token"`
}

// Login godoc
// @Summary User login
// @Description Log in a user using email/phone and password.
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body LoginReq true "Login details"
// @Success 200 {object} LoginRes "Login successful, returns JWT token"
// @Failure 400 {object} Response "{"hasError": "true", "message": "Bad request"}"
// @Failure 400 {object} Response "{"hasError": "true", "message": "Invalid email format"}"
// @Failure 400 {object} Response "{"hasError": "true", "message": "Invalid phone number format"}"
// @Failure 400 {object} Response "{"hasError": "true", "message": "Please input your password!"}"
// @Failure 400 {object} Response "{"hasError": "true", "message": "This user not valid!"}"
// @Failure 400 {object} Response "{"message": "Bad request"}"
// @Security ApiKeyAuth
// @Router /login [post]
func (auth *Auth) Login(ctx echo.Context) error {

	json_map := make(map[string]interface{})
	var loginReq LoginReq

	err := json.NewDecoder(ctx.Request().Body).Decode(&json_map)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Bad request"})
	}

	if json_map["email"] != nil {
		loginReq.Email = json_map["email"].(string)
		emailError := utils.IsValidEmail(loginReq.Email)
		if emailError != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"hasError": "true", "message": emailError.Error()})
		}
	}

	if json_map["phoneNumber"] != nil {
		loginReq.PhoneNumber = json_map["phoneNumber"].(string)
		phoneNumberError := utils.IsValidPhoneNumber(json_map["phoneNumber"].(string))
		if phoneNumberError != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"hasError": "true", "message": phoneNumberError.Error()})
		}
	}

	if json_map["password"] != nil {
		loginReq.Password = json_map["password"].(string)
		if len(loginReq.Password) == 0 {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"hasError": "true", "message": "Please input your password!"})
		}
	} else {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"hasError": "true", "message": "Please input your password!"})
	}

	user, errValidationUser := utils.GetUser(auth.DB, loginReq.Email, loginReq.PhoneNumber, loginReq.Password)

	if errValidationUser != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"hasError": "true", "message": "This user not valid!"})
	}

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Bad request"})
	}

	token, err := utils.CreateToken(auth.JWT, int(user.ID))

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Bad request"})
	}

	jsonToken := new(LoginRes)
	jsonToken.Token = token

	return ctx.JSON(http.StatusOK, jsonToken)
}
