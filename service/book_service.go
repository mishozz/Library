package service

import (
	"github.com/mishozz/Library/entities"
	"github.com/mishozz/Library/repositories"
)

type BookService interface {
	Save(entities.Book) error
	FindAll() []entities.Book
}

type bookService struct {
	repository repositories.BookRepository
}

func NewBookService(repo repositories.BookRepository) BookService {
	return &bookService{
		repository: repo,
	}
}

func (s *bookService) Save(book entities.Book) error {
	s.repository.Save(book)
	return nil
}

func (s *bookService) FindAll() []entities.Book {
	return s.repository.FindAll()
}
