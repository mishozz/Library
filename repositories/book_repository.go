package repositories

import (
	"github.com/mishozz/Library/config"
	"github.com/mishozz/Library/entities"
	"gorm.io/gorm"
)

type BookRepository interface {
	Save(book entities.Book)
	Delete(isbn string) error
	FindAll() []entities.Book
	Find(isbn string) entities.Book
	UpdateUnits(book entities.Book)
	IsBookTaken(isbn string) bool
}

type BookRepositoryImpl struct {
	connection *gorm.DB
}

func NewBookRepository(db config.Database) *BookRepositoryImpl {
	return &BookRepositoryImpl{
		connection: db.Connection,
	}
}

func (b *BookRepositoryImpl) Save(book entities.Book) {
	b.connection.Create(&book)
}

func (b *BookRepositoryImpl) Delete(isbn string) error {
	var book entities.Book
	b.connection.Where("Isbn = ?", isbn).First(&book)

	if err := b.connection.Model(&book).Association("UserReturned").Delete(&book); err != nil {
		return err
	}
	if err := b.connection.Unscoped().Delete(&book).Error; err != nil {
		return err
	}
	return nil
}

func (b *BookRepositoryImpl) IsBookTaken(isbn string) bool {
	var book entities.Book
	b.connection.Where("Isbn = ?", isbn).First(&book)

	var users []entities.User
	b.connection.Model(&book).Association("UserTaken").Find(&users)
	//b.connection.Unscoped().Table("user_taken").Where("book_id = ?", book.Model.ID).Find(&users)
	if len(users) != 0 {
		return true
	}

	return false
}

func (b *BookRepositoryImpl) FindAll() []entities.Book {
	var books []entities.Book
	b.connection.Set("gorm:auto_preload", true).Find(&books)
	return books
}

func (b *BookRepositoryImpl) Find(isbn string) entities.Book {
	var book entities.Book
	b.connection.Where("Isbn = ?", isbn).First(&book)
	return book
}
func (b *BookRepositoryImpl) UpdateUnits(book entities.Book) {
	b.connection.Model(&book).Update("available_units", book.AvailableUnits)
}
