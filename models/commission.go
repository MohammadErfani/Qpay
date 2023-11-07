package models

import (
	"gorm.io/gorm"
	"time"
)

type Commission struct {
	gorm.Model
	ComID              uint `gorm:"primaryKey"`
	AmountPerTrans     float64
	PercentagePerTrans float64
	Status             uint8
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          time.Time
}
