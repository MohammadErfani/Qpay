package handlers

import (
	"Qpay/database"
	"Qpay/models"
	"Qpay/services/auth"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
)

type RegsiterRequest struct {
	Name        string `json:"email" binding:"required"`
	Username    string `json:"email" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"email" binding:"required"`
	PhoneNumber string `json:"email" binding:"required"`
	Identity    string `json:"email" binding:"required"`
	Address     string `json:"email" binding:"required"`
	IsUser      bool   `json:"email" binding:"required"` // 0 , 1 , 2
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
	phoneRegex := regexp.MustCompile(`^\d{9}$`)
	if !phoneRegex.MatchString(user.PhoneNumber) {
		return errors.New("invalid phone number format")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		return errors.New("invalid email format")
	}
	identityRegex := regexp.MustCompile(`^\d{10}$`)
	if !identityRegex.MatchString(user.Identity) {
		return errors.New("identity number must be 10 digits")
	}
	if len(user.Password) < 4 {
		return errors.New("password must be at least 4 characters long")
	}
	return nil
}
