package entities

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Isbn           string `json:"Isbn" binding:"required" gorm:"type:varchar(32);UNIQUE"`
	Title          string `json:"Title" binding:"required" gorm:"type:varchar(256)"`
	Author         string `json:"Author" binding:"required" gorm:"type:varchar(100)"`
	AvailableUnits uint   `json:"AvailableUnits" binding:"required"`
}