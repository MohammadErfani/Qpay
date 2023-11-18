package routes

import (
	"Qpay/server/handlers"
	"github.com/labstack/echo/v4"
)

func GatewayGroup(ctx *echo.Group) {
	ctx.GET("/gateway", handlers.ListAllGateways)     // List all gateways
	ctx.GET("/gateway/:id", handlers.FindGateway)     // find a gateways
	ctx.PUT("/gateway/:id", handlers.DeleteCard)      // update gateway for a user
	ctx.POST("/gateway", handlers.RegisterNewGateway) // register gateway for a user
}
