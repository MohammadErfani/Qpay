package routes

import (
	"Qpay/server/handlers"
	"github.com/labstack/echo/v4"
)

func GatewayGroup(ctx *echo.Group) {
	h := &handlers.Handler{
		UserID: 1,
	}
	ctx.GET("/gateway", h.ListAllGateways)     // List all gateways
	ctx.GET("/gateway/:id", h.FindGateway)     // find a gateways
	ctx.PATCH("/gateway/:id", h.UpdateGateway) // update gateway for a user
	ctx.POST("/gateway", h.RegisterNewGateway) // register gateway for a user
}
