package server

import (
	"Qpay/config"
	"Qpay/server/routes"
	"fmt"
	"log"
	"sync"

	"github.com/labstack/echo/v4"
)

var (
	once sync.Once
	srv  *Server
)

type Server struct {
	router *echo.Echo
}

func Instance() *Server {
	once.Do(func() {
		srv = &Server{
			router: routes.InitRoutesV1(),
		}
	})
	return srv
}

func StartServer(cfg *config.Config, srv *Server) {
	log.Fatal(srv.router.Start(fmt.Sprintf(":%d", cfg.Server.Port)))
}
