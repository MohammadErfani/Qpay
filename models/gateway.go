package models

import (
	"gorm.io/gorm"
)

type Gateway struct {
	gorm.Model
	UserID        uint
	CommissionID  uint
	BankAccountID uint
	User          User        `gorm:"foreignKey:UserID"`
	Commission    Commission  `gorm:"foreignKey:CommissionID"`
	BankAccount   BankAccount `gorm:"foreignKey:BankAccountID"`
	Transactions  []Transaction
	Name          string `gorm:"type:varchar(50)"`
	Logo          string `gorm:"type:varchar(50)"`
	Route         string `gorm:"type:varchar(50)"`
	Status        uint8
	Type          uint8
}

const (
	StatusGatewayActive = iota
	StatusGatewayInActive
	StatusGatewayUnapproved
	StatusGatewayDraft
)

const (
	PersonalTypeGateway = iota
	BusinessTypeGateway
)

func SetGatewayType(isPersonal bool) uint8 {
	if isPersonal {
		return PersonalTypeGateway
	}
	return BusinessTypeGateway
}
