package routes

import (
	"Qpay/database"
	"Qpay/server/handlers"
	echo "github.com/labstack/echo/v4"
)

func BankAccountGroup(bc *echo.Group) {
	h := &handlers.Handler{
		DB: database.DB(),
	}
	bc.GET("/bankaccount", h.ListAllBankAccounts)      // List all cards
	bc.GET("/bankaccount/:id", h.FindBankAccount)      // find a card
	bc.POST("/bankaccount", h.RegisterNewBankAccount)  // register card for a user
	bc.DELETE("/bankaccount/:id", h.DeleteBankAccount) // delete card for a user
}

/*
HTTP Verb		Path (URL)			Action (Method)		Route Name
GET				/sharks				index				sharks.index
GET				/sharks/create		create				sharks.create
POST			/sharks				store				sharks.store
GET				/sharks/{id}		show				sharks.show
GET				/sharks/{id}/edit	edit				sharks.edit
PUT/PATCH		/sharks/{id}		update				sharks.update
DELETE			/sharks/{id}		destroy				sharks.destroy
*/
