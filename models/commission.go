package models

import (
	"gorm.io/gorm"
)

type Commission struct {
	gorm.Model
	ComID              uint `gorm:"primaryKey"`
	AmountPerTrans     float64
	PercentagePerTrans float64
	Status             uint8
	Gateways           []Gateway
}
