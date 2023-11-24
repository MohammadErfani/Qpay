package routes

import (
	"Qpay/database"
	"Qpay/server/handlers"
	"github.com/labstack/echo/v4"
)

func GatewayGroup(ctx *echo.Group) {
	gh := &handlers.GatewayHandler{
		UserID: 1,
		DB:     database.DB(),
	}
	ctx.GET("/gateway", gh.ListAllGateways) // List all gateways
	ctx.GET("/gateway/:id", gh.FindGateway) // find a gateways
	ctx.POST("/gateway", gh.RegisterNewGateway)
	ctx.POST("/gateway/update", gh.ChangeGetawayBankAccount)
}
