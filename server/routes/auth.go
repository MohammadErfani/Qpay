package routes

import (
	"Qpay/config"
	"Qpay/server/handlers"
	"Qpay/server/middlewares"

	echo "github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type User struct {
	email       string `json:"email"`
	phoneNumber string `json:"phoneNumber"`
	password    string `json:"password"`
}

type Auth struct {
	DB  *gorm.DB
	JWT *config.JWT
}

func AuthGroup(authG *echo.Group, db *gorm.DB, cfg *config.Config) {

	auth := &handlers.Auth{
		DB:  db,
		JWT: &cfg.JWT,
	}
	// authG.POST("/login", auth.Login)
	authG.POST("/login", middlewares.AuthMiddleware(auth.Login))
}
