package models

import (
	"gorm.io/gorm"
)

type Commission struct {
	gorm.Model
	AmountPerTrans     float64
	PercentagePerTrans float64
	Status             uint8
	Gateways           []Gateway
}

const (
	CommIsActive = iota
	CommIsInactive
)
