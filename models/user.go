package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	UserID       uint   `gorm:"primaryKey"`
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
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}
