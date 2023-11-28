package handlers

import (
	"Qpay/server/middlewares"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	UserID uint
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (h *Handler) SetUserID(ctx echo.Context) {
	userID := ctx.Get(middlewares.UserIdContextField).(int)
	h.UserID = uint(userID)
}
