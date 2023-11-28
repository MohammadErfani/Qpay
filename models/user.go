package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"type:varchar(50)"`
	Username     string `gorm:"type:varchar(50);unique"`
	Email        string `gorm:"type:varchar(50);unique"`
	Password     string `gorm:"type:varchar(128)"`
	PhoneNumber  string `gorm:"type:varchar(15);unique"`
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

func SetRole(isCompany bool) uint8 {
	if isCompany {
		return IsCompany
	}
	return IsNaturalPerson
}
