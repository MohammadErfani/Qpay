package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	GatewayID        uint
	Gateway          Gateway `gorm:"foreignKey:GatewayID"`
	PaymentAmount    float64
	CommissionAmount float64
	PurchaserCard    string
	CardMonth        int
	CardYear         int
	PhoneNumber      string
	//کد رهگیری
	Status               uint8
	OwnerBankAccount     string `gorm:"type:varchar(50)"`
	PurchaserBankAccount string `gorm:"type:varchar(50)"`
}

const (
	NotPaid = iota
	NotSuccessfully
	IssueOccurred
	Blocked
	Refund
	Cancelled
	ReturnToGateway
	AwaitingConfirmation
	Confirmed
	newConst
)
