package handlers

import (
	"Qpay/database"
	"Qpay/models"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Card struct {
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
	sheba := ctx.QueryParam("sheba")
	card := models.BankAccount{Sheba: sheba}
	if err := ctx.Bind(&card); err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, card)
}

func DeleteCard(ctx echo.Context) error {
	return nil
}

/**


func(c echo.Context) error {
		db := database.DB()
		//card := models.BankAccount{
		//	Sheba:  "IR0696000000010324200001",
		//	Status: 1,
		//	UserID: 1,
		//}
		//_, err := bankaccount.CreateCard(db, card)
		//if err != nil {
		//	return c.JSON(http.StatusBadRequest, "Internal server error in create user ")
		//}
		//account := &models.BankAccount{
		//	UserID:       1,
		//	BankID:       1,
		//	Status:       0,
		//	AccountOwner: "amin",
		//	Sheba:        "IR0696000000010324200001",
		//}
		//acc := &models.BankAccount{
		//	UserID:       1,
		//	BankID:       1,
		//	User:         models.User{Username: "mohammadErfani"},
		//	Bank:         models.Bank{Name: "بانک ملی ایران"},
		//	Status:       0,
		//	AccountOwner: "محمد عرفانی",
		//	Sheba:        "IR0696000000010324200001",
		//}
		//db.Save(acc)
		var card01 models.BankAccount
		db.First(&card01)
		return json.NewEncoder(c.Response()).Encode(card01)
	}



*/
