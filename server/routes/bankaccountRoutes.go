package routes

import (
	"Qpay/database"
	"Qpay/models"
	"encoding/json"
	echo "github.com/labstack/echo/v4"
)

//var cfgFile string = "config.yaml"

func BankAccountRoutes(bc *echo.Group) {
	bc.GET("/bankaccount/readall", func(c echo.Context) error {
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

		var card01 models.BankAccount
		db.First(&card01, 1)
		//db.Save(account)
		return json.NewEncoder(c.Response()).Encode(card01)
	})
}
