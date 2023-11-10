package models

import (
	"gorm.io/gorm"
	"time"
)

type Bank struct {
	gorm.Model
	BankID       uint   `gorm:"primaryKey"`
	Name         string `gorm:"type:varchar(50)"`
	Logo         string
	BankAccounts []BankAccount
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}
