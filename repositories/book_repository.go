package repositories

import (
	"github.com/mishozz/Library/config"
	"github.com/mishozz/Library/entities"
	"gorm.io/gorm"
)

type BookRepository interface {
	Save(book entities.Book) error
	Delete(isbn string) error
	FindAll() ([]entities.Book, error)
	Find(isbn string) (entities.Book, error)
	UpdateUnits(book entities.Book) error
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

func (b *BookRepositoryImpl) Save(book entities.Book) error {
	err := b.connection.Create(&book).Error
	if err != nil {
		return err
	}
	return nil
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
	if len(users) != 0 {
		return true
	}

	return false
}

func (b *BookRepositoryImpl) FindAll() ([]entities.Book, error) {
	var books []entities.Book
	err := b.connection.Set("gorm:auto_preload", true).Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (b *BookRepositoryImpl) Find(isbn string) (entities.Book, error) {
	var book entities.Book
	err := b.connection.Where("Isbn = ?", isbn).First(&book).Error
	if err != nil {
		return book, err
	}
	return book, nil
}
func (b *BookRepositoryImpl) UpdateUnits(book entities.Book) error {
	err := b.connection.Model(&book).Update("available_units", book.AvailableUnits).Error
	if err != nil {
		return err
	}
	return nil
}
