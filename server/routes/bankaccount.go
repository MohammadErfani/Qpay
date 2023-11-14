package routes

import (
	"Qpay/server/handlers"
	echo "github.com/labstack/echo/v4"
)

func BankAccount(bc *echo.Group) {
	bc.GET("/bankaccount", handlers.ListAllCards)             // List all cards
	bc.GET("/bankaccount/:id", handlers.FindCard)             // find a card
	bc.POST("/bankaccount", handlers.RegisterNewCard)         // register card for a user
	bc.DELETE("/bankaccount/delete/:id", handlers.DeleteCard) // delete card for a user
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
