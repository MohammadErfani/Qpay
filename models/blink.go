package models

import "gorm.io/gorm"

type Blinks struct {
	gorm.Model
	UserID uint
	Price  uint
	User   User   `gorm:"foreignKey:UserID"`
	Name   string `gorm:"type:varchar(50)"`
	Status uint8
}

const (
	LinkIsAvailable = iota
	LinkIsSold
	LinkIsUnAvailable
)
