package service

import (
	"reflect"

	"github.com/mishozz/Library/entities"
	"github.com/mishozz/Library/repositories"
)

type BookService interface {
	Save(entities.Book)
	FindAll() []entities.Book
	FindByIsbn(isbn string) entities.Book
	BookExists(isbn string) bool
}

type bookService struct {
	repository repositories.BookRepository
}

func NewBookService(repo repositories.BookRepository) BookService {
	return &bookService{
		repository: repo,
	}
}

func (s *bookService) Save(book entities.Book) {
	s.repository.Save(book)
}

func (s *bookService) FindAll() []entities.Book {
	return s.repository.FindAll()
}

func (s *bookService) FindByIsbn(isbn string) entities.Book {
	return s.repository.Find(isbn)
}

func (s *bookService) BookExists(isbn string) bool {
	book := s.repository.Find(isbn)
	return !reflect.ValueOf(book).IsZero()
}
