package handlers

import (
	"Qpay/services"
	"Qpay/utils"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type AdminRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminHandler struct {
	DB *gorm.DB
}

type AdminResponse struct {
}

func (aH *AdminHandler) AdminCreate(ctx echo.Context) error {
	var req AdminRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	if err := ValidateAdmin(&req); err != nil {
		return ctx.JSON(http.StatusForbidden, err.Error())
	}
	if err := ValidateUniqueAdmin(aH.DB, &req); err != nil {
		return ctx.JSON(http.StatusConflict, err.Error())
	}
	_, err := services.CreateAdmin(
		aH.DB,
		req.Name,
		req.Username,
		req.Email,
		req.Password,
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server error in create admin")
	}
	return ctx.JSON(http.StatusCreated, RegisterResponse{Status: "success", Message: "Admin created successfully"})

}

func ValidateAdmin(admin *AdminRequest) error {
	requiredFields := map[string]string{
		"name":     admin.Name,
		"email":    admin.Email,
		"password": admin.Password,
		"username": admin.Username,
	}
	if err := utils.IsRequired(requiredFields); err != nil {
		return err
	}
	if err := utils.IsValidEmail(admin.Email); err != nil {
		return err
	}
	return nil
}
func ValidateUniqueAdmin(db *gorm.DB, admin *AdminRequest) error {
	uniqueFields := map[string]string{
		"username": admin.Username,
		"email":    admin.Email,
	}
	for fieldName, fieldValue := range uniqueFields {
		if _, err := services.GetUser(db, fieldName, fieldValue); err == nil {
			return errors.New(fmt.Sprintf("%s already exist", fieldName))
		}
	}
	return nil
}
