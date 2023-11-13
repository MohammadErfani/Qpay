package handlers

import (
	"Qpay/database"
	"Qpay/models"
	"Qpay/services/auth"
	"Qpay/utils"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type RegsiterRequest struct {
	Name        string `json:"name" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Identity    string `json:"identity" binding:"required"`
	Address     string `json:"address" binding:"required"`
	IsUser      bool   `json:"isUser" binding:"required"` // 0 , 1 , 2
}
type RegisterResponse struct {
	Status  string
	Message string
}

func CreateUser(ctx echo.Context) error {
	db := database.DB()
	var req RegsiterRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	if err := ValidateUser(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	existUser, _ := auth.GetUserByEmail(db, req.Email)
	if existUser != nil {
		return ctx.JSON(http.StatusBadRequest, "User Email exist")
	}
	existUser, _ = auth.GetUserByUsername(db, req.Username)
	if existUser != nil {
		return ctx.JSON(http.StatusBadRequest, "User name exist")
	}
	user := models.User{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Identity:    req.Identity,
		Address:     req.Address,
	}
	_, err := auth.CreateUser(db, user)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Internal server error in create user ")
	}
	return ctx.JSON(http.StatusCreated, RegisterResponse{Status: "success", Message: "User created successfully"})

}
func ValidateUser(user *RegsiterRequest) error {
	if user.Name == "" {
		return errors.New("name is required")
	}
	if user.Username == "" {
		return errors.New("username is required")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Password == "" {
		return errors.New("phone number is required")
	}
	if utils.IsValidEmail(user.Email) {
		return errors.New("email is not in correct format")
	}
	if utils.IsValidPhoneNumber(user.PhoneNumber) {
		return errors.New("phone number is not in correct format")
	}
	if utils.IsValidNationalCode(user.Identity) {
		return errors.New("identity is not in correct format")
	}
	return nil
}
