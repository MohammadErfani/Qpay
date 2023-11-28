package routes

import (
	"Qpay/database"
	"Qpay/server/handlers"
	"github.com/labstack/echo/v4"
)

func AdminGroup(adminG *echo.Group) {
	h := handlers.Handler{
		DB: database.DB(),
	}
	adminG.POST("/admin/register", h.AdminCreate)
	adminG.GET("/admin/commission", h.AdminListAllCommission)
	adminG.GET("/admin/commission/:id", h.AdminGetCommission)
	adminG.POST("/admin/commission", h.AdminCreateCommission)
	adminG.GET("/admin/user", h.AdminListUsers)
	adminG.GET("/admin/user/:id", h.AdminGetUser)
	adminG.PATCH("/admin/user/:id", h.AdminUpdateUser)
	adminG.GET("/admin/gateway", h.AdminListAllGateways)
	adminG.GET("/admin/gateway/:id", h.AdminGetGateway)
	adminG.PATCH("/admin/gateway/:id", h.AdminUpdateGateway)
	adminG.GET("/admin/transaction", h.AdminListTransactions)
	adminG.GET("/admin/transaction/:id", h.AdminGetTransaction)
	adminG.PATCH("/admin/transaction/:id", h.AdminUpdateTransaction)
}
