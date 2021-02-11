package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model    `json:"-"`
	Email         string `json:"Email" binding:"required" gorm:"type:varchar(100);UNIQUE"`
	Password      string `json:"Password,omitempty"`
	Role          string `gorm:"size:255;not null;" json:"-"`
	TakenBooks    []Book `json:"Taken_books" gorm:"many2many:user_taken;"`
	ReturnedBooks []Book `json:"Returned_books" gorm:"many2many:user_returned;"`
}
