package cmd

import (
	"Qpay/config"
	"Qpay/database"
	"Qpay/server"
)

func Execute() {
	cfg := config.InitConfig("config.yaml")
	db := database.NewPostgres(cfg)
	_ = db
	server.StartServer(cfg, server.Instance())
}
