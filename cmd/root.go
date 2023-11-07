package cmd

import (
	"Qpay/config"
	"Qpay/database"
)

func Execute() {
	cfg := config.InitConfig("config.yaml")
	db := database.NewPostgres(cfg)
	_ = db
}
