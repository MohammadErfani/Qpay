package handlers

import (
	"Qpay/database"
	"Qpay/models"
	"Qpay/services"
	"Qpay/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type BankAccountRequest struct {
	Sheba string `json:"sheba" xml:"sheba" form:"sheba" query:"sheba"`
}

func ListAllCards(ctx echo.Context) error {
	db := database.DB()
	var cards []models.BankAccount
	var userID uint = 1 //user id ro bayad tashkhis bedim o inja vared konim.
	result := db.Where(&models.BankAccount{UserID: userID}).Order("created_at desc").Find(&cards)
	if result.RowsAffected == 0 {
		return ctx.JSON(http.StatusBadRequest, "You Aren't Add Any Card. Please add a card!")
	}
	return json.NewEncoder(ctx.Response()).Encode(&cards)
}
func FindCard(ctx echo.Context) error {
	db := database.DB()
	var card models.BankAccount
	CardID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	result := db.First(&card, CardID)
	if result.RowsAffected == 0 {
		return ctx.JSON(http.StatusNotFound, "You're Card is not exist!")
	}
	return json.NewEncoder(ctx.Response()).Encode(&card)
}

func RegisterNewCard(ctx echo.Context) error {
	db := database.DB()
	var req BankAccountRequest
	//var userID uint = 1 //user id ro bayad tashkhis bedim o inja vared konim.
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	//sheba := ctx.QueryParam("sheba")
	err := utils.CheckSheba(req.Sheba)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if err := ValidateUniqueBankAccount(db, &req); err != nil {
		return ctx.JSON(http.StatusConflict, err.Error())
	}
	_, err = services.CreateBankAccount(db, req.Sheba)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server error in create bank account")
	}
	//bankID, Owner := utils.GetOwnerSheba(sheba)
	//card := models.BankAccount{Sheba: sheba, BankID: bankID, AccountOwner: Owner, UserID: userID}
	//db.Save(&card)
	return ctx.JSON(http.StatusOK, "You're card is successfully registered!")
}

func DeleteCard(ctx echo.Context) error {
	db := database.DB()
	id := ctx.Param("id")
	fmt.Println(id)
	err := db.Delete(&models.BankAccount{}, "id=?", id).Error
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}
	return ctx.JSON(http.StatusOK, "You're Card is successfully deleted!")
}

func ValidateUniqueBankAccount(db *gorm.DB, bankAccount *BankAccountRequest) error {
	if _, err := services.GetBankAccount(db, "sheba", bankAccount.Sheba); err == nil {
		return errors.New("sheba already exist")
	}
	return nil
}
