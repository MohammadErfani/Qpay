package models

import (
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	gorm.Model
	TransID              uint `gorm:"primaryKey"`
	GatewayID            uint
	Gateway              Gateway `gorm:"foreignKey:GatewayID"`
	PaymentAmount        float64
	Status               uint8
	OwnerBankAccount     string `gorm:"type:varchar(50)"`
	PurchaserBankAccount string `gorm:"type:varchar(50)"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            time.Time
}
