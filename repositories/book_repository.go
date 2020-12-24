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
	Find(isbn string)
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

}

func (b *BookRepositoryImpl) Delete(book entities.Book) {

}

func (b *BookRepositoryImpl) FindAll() []entities.Book {
	return nil
}

func (b *BookRepositoryImpl) Find(isbn string) {

}
