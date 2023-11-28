package handlers

import (
	"Qpay/models"
	"Qpay/services"
	"Qpay/utils"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type RegisterRequest struct {
	Name        string `json:"name" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Identity    string `json:"identity" binding:"required"`
	Address     string `json:"address" binding:"required"`
	IsCompany   bool   `json:"is_company" binding:"required"`
}

// @Summary Create a new user
// @Description Create a new user with the provided information.
// @Tags users
// @Accept json
// @Produce json
// @Param registerRequest body RegisterRequest true "User registration details"
// @Success 201 {object} Response "{"status": "success", "message": "User created successfully"}"
// @Failure 400 {object} Response "{"status": "error", "message": "Bind Error"}"
// @Failure 403 {object} Response "{"status": "error", "message": "Validation Error"}"
// @Failure 409 {object} Response "{"status": "error", "message": "Conflict - User already exists"}"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error in creating user"}"
// @Router /register [post]
func (h *Handler) CreateUser(ctx echo.Context) error {
	var req RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Bind Error",
		})
	}
	if err := ValidateUser(&req); err != nil {
		return ctx.JSON(http.StatusForbidden, Response{
			Status:  "error",
			Message: err.Error(),
		})
	}
	if err := ValidateUserUnique(h.DB, &req); err != nil {
		return ctx.JSON(http.StatusConflict, Response{
			Status:  "error",
			Message: err.Error(),
		})
	}
	user := models.User{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Identity:    req.Identity,
		Address:     req.Address,
		Role:        models.SetRole(req.IsCompany),
	}
	_, err := services.CreateUser(h.DB, user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Internal server error in creating user",
		})
	}
	return ctx.JSON(http.StatusCreated, Response{Status: "success", Message: "User created successfully"})

}
func ValidateUser(user *RegisterRequest) error {

	requiredFields := map[string]string{
		"name":         user.Name,
		"email":        user.Email,
		"password":     user.Password,
		"phone_number": user.PhoneNumber,
		"identity":     user.Identity,
	}
	if err := utils.IsRequired(requiredFields); err != nil {
		return err
	}
	if err := utils.IsValidEmail(user.Email); err != nil {
		return err
	}
	if err := utils.IsValidPhoneNumber(user.PhoneNumber); err != nil {
		return err
	}
	if err := utils.IsValidNationalCode(user.Identity); err != nil {
		return err
	}
	return nil
}

func ValidateUserUnique(db *gorm.DB, user *RegisterRequest) error {

	uniqueFields := map[string]string{
		"username":     user.Username,
		"email":        user.Email,
		"phone_number": user.PhoneNumber,
		"identity":     user.Identity,
	}
	for fieldName, fieldValue := range uniqueFields {
		if _, err := services.GetUser(db, fieldName, fieldValue); err == nil {
			return errors.New(fmt.Sprintf("%s already exist", fieldName))
		}
	}
	return nil
}
