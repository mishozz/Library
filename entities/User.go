package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email         string `json:"Email" binding:"required" gorm:"type:varchar(100);UNIQUE"`
	TakenBooks    []Book `json:"Taken_books" gorm:"many2many:user_taken;"`
	ReturnedBooks []Book `json:"Returned_books" gorm:"many2many:user_returned;"`
}
