package routes

import (
	"Qpay/database"
	"Qpay/server/handlers"
	"github.com/labstack/echo/v4"
)

func AdminGroup(adminG *echo.Group) {
	aH := handlers.AdminHandler{
		DB: database.DB(),
	}
	adminG.POST("/admin/register", aH.AdminCreate)
	adminG.GET("/admin/commission", aH.AdminListAllCommission)
	adminG.GET("/admin/commission/:id", aH.AdminGetCommission)
	adminG.POST("/admin/commission", aH.AdminCreateCommission)
	adminG.GET("/admin/user", aH.AdminListUsers)
	adminG.GET("/admin/user/:id", aH.AdminGetUser)
	adminG.PATCH("/admin/user/:id", aH.AdminUpdateUser)
	adminG.GET("/admin/gateway", aH.AdminListAllGateways)
	adminG.GET("/admin/gateway/:id", aH.AdminGetGateway)
	adminG.PATCH("/admin/gateway/:id", aH.AdminUpdateGateway)
	adminG.GET("/admin/transaction", aH.AdminListTransactions)
	adminG.GET("/admin/transaction/:id", aH.AdminGetTransaction)
	adminG.PATCH("/admin/transaction/:id", aH.AdminUpdateTransaction)
}
