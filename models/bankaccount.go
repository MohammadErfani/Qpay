package models

import (
	"gorm.io/gorm"
	"time"
)

type BankAccount struct {
	gorm.Model
	BankAccountID uint `gorm:"primaryKey"`
	UserID        uint
	BankID        uint
	User          User `gorm:"foreignKey:UserID"`
	Bank          Bank `gorm:"foreignKey:BankID"`
	Gateways      []Gateway
	Status        uint8
	AccountOwner  string `gorm:"type:varchar(50)"`
	Sheba         string `gorm:"type:varchar(50);unique"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}
