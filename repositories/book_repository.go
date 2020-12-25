package repositories

import (
	"github.com/mishozz/Library/config"
	"github.com/mishozz/Library/entities"
	"gorm.io/gorm"
)

type BookRepository interface {
	Save(book entities.Book)
	Delete(book entities.Book)
	FindAll() []entities.Book
	Find(isbn string) entities.Book
}

type BookRepositoryImpl struct {
	connection *gorm.DB
}

func NewBookRepository(db config.Database) BookRepository {
	return &BookRepositoryImpl{
		connection: db.Connection,
	}
}

func (b *BookRepositoryImpl) Save(book entities.Book) {
	b.connection.Create(&book)
}

func (b *BookRepositoryImpl) Delete(book entities.Book) {

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
