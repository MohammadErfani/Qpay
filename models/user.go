package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"type:varchar(50)"`
	Username     string `gorm:"type:varchar(50);unique"`
	Email        string `gorm:"type:varchar(50);unique"`
	Password     string `gorm:"type:varchar(128);unique"`
	PhoneNumber  string `gorm:"type:varchar(15)"`
	Identity     string `gorm:"type:varchar(50);unique"`
	Address      string `gorm:"type:text;unique"`
	Role         uint8
	Gateways     []Gateway
	BankAccounts []BankAccount
}

const (
	IsNaturalPerson = iota
	IsCompany
	IsAdmin
)
