package handlers

import (
	"Qpay/database"
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
	IsCompany   bool   `json:"is_company" binding:"required"` // 0 , 1 , 2
}
type RegisterResponse struct {
	Status  string
	Message string
}

// @Summary Create a new user
// @tags Register
// @Description Create a new user
// @ID create-user
// @Accept  json
// @Produce  json
// @Param body body RegisterRequest true "User registration request"
// @Success 201 {object} RegisterResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 403 {string} string "validation error"
// @Failure 409 {string} string "duplicate email, username, phone or identity"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/register [post]
func CreateUser(ctx echo.Context) error {
	db := database.DB()
	var req RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	if err := ValidateUser(&req); err != nil {
		return ctx.JSON(http.StatusForbidden, err.Error())
	}
	if err := ValidateUserUnique(db, &req); err != nil {
		return ctx.JSON(http.StatusConflict, err.Error())
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
	_, err := services.CreateUser(db, user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server error in create user ")
	}
	return ctx.JSON(http.StatusCreated, RegisterResponse{Status: "success", Message: "User created successfully"})

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
