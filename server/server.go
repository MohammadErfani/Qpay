package server

import (
	"Qpay/config"
	"Qpay/server/routes"
	"fmt"
	"log"
	"sync"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	once sync.Once
	srv  *Server
)

type Server struct {
	router *echo.Echo
}

func Instance(db *gorm.DB, cfg *config.Config) *Server {
	once.Do(func() {
		srv = &Server{
			router: routes.InitRoutesV1(db, cfg),
		}
	})
	return srv
}

func StartServer(cfg *config.Config, srv *Server) {
	log.Fatal(srv.router.Start(fmt.Sprintf(":%d", cfg.Server.Port)))
}
