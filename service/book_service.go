package service

import (
	"github.com/mishozz/Library/entities"
	"github.com/mishozz/Library/repositories"
)

type BookService interface {
	Save(entities.Book) error
	FindAll() ([]entities.Book, error)
	FindByIsbn(isbn string) (entities.Book, error)
	Delete(isbn string) error
	IsBookTaken(isbn string) bool
}

type bookService struct {
	repository repositories.BookRepository
}

func NewBookService(repo repositories.BookRepository) *bookService {
	return &bookService{
		repository: repo,
	}
}

func (s *bookService) Save(book entities.Book) error {
	return s.repository.Save(book)
}

func (s *bookService) FindAll() ([]entities.Book, error) {
	return s.repository.FindAll()
}

func (s *bookService) FindByIsbn(isbn string) (entities.Book, error) {
	return s.repository.Find(isbn)
}

func (s *bookService) Delete(isbn string) error {
	return s.repository.Delete(isbn)
}

func (s *bookService) IsBookTaken(isbn string) bool {
	return s.repository.IsBookTaken(isbn)
}
