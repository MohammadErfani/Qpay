package handlers

import (
	"Qpay/config"
	"Qpay/database"
	"Qpay/models"
	"gorm.io/gorm"
	"log"
)

type Auth struct {
	DB *gorm.DB
}

func addCard(configPath string) {
	cfg := config.InitConfig(configPath)
	db := database.NewPostgres(cfg)
	card := models.BankAccount{
		Sheba:  "Mohammad Erfani",
		UserID: 1,
		Status: 0,
	}
	err := db.FirstOrCreate(&card, models.BankAccount{Sheba: card.Sheba, UserID: card.UserID, Status: card.Status}).Error
	if err != nil {
		log.Fatal(err)
	}
}
