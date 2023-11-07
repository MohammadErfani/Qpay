package database

import (
	"Qpay/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/url"
	"sync"
	"time"
)

var instance *gorm.DB
var once sync.Once

func NewPostgres(cfg *config.Config) *gorm.DB {
	once.Do(func() {
		tehranTimezone, _ := time.LoadLocation("Asia/Tehran")
		// Connection configuration
		dsn := &url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword(cfg.Database.Username, cfg.Database.Password),
			Host:     fmt.Sprintf("localhost:%d", cfg.Database.Port),
			Path:     cfg.Database.DB,
			RawQuery: "sslmode=disable&timezone=" + tehranTimezone.String(),
		}

		// Convert URL to connection string
		connStr := dsn.String()

		db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to the database:", err)
		}

		fmt.Println("Successfully connected to the database!")

		instance = db
	})

	return instance

}
