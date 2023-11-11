package models

import (
	"gorm.io/gorm"
)

type Bank struct {
	gorm.Model
	Name         string `gorm:"type:varchar(50)"`
	Logo         string
	BankAccounts []BankAccount
}
