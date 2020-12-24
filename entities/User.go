package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email         string `json:"Email" binding:"required" gorm:"type:varchar(100);UNIQUE"`
	TakenBooks    []Book `gorm:"many2many:user_taken;"`
	ReturnedBooks []Book `gorm:"many2many:user_returned;"`
}
