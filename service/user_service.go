package service

import (
	"reflect"

	"github.com/mishozz/Library/entities"
	"github.com/mishozz/Library/repositories"
)

type UserService interface {
	FindByEmail(email string) entities.User
	FindAll() []entities.User
	UserExists(email string) bool
	TakeBook(user entities.User, book entities.Book)
	ReturnBook(user entities.User, book entities.Book)
	IsBookTakenByUser(email string, isbn string) bool
}

type userService struct {
	userRepository repositories.UserRepository
	bookRepository repositories.BookRepository
}

func NewUserService(userRepository repositories.UserRepository, bookRepository repositories.BookRepository) *userService {
	return &userService{
		userRepository: userRepository,
		bookRepository: bookRepository,
	}
}

func (s *userService) FindByEmail(email string) entities.User {
	return s.userRepository.FindByEmail(email)
}

func (s *userService) FindAll() []entities.User {
	return s.userRepository.FindAll()
}

func (s *userService) UserExists(email string) bool {
	user := s.userRepository.FindByEmail(email)
	return !reflect.ValueOf(user).IsZero()
}

func (s *userService) TakeBook(user entities.User, book entities.Book) {
	book.AvailableUnits = book.AvailableUnits - 1
	user.TakenBooks = append(user.TakenBooks, book)

	s.bookRepository.UpdateUnits(book)
	s.userRepository.UpdateTakenBooks(user, user.TakenBooks)

	if contains(user.ReturnedBooks, book) {
		user.ReturnedBooks = remove(user.ReturnedBooks, book)
		s.userRepository.UpdateReturnedBooks(user, user.ReturnedBooks)
	}
}

func (s *userService) ReturnBook(user entities.User, book entities.Book) {
	book.AvailableUnits = book.AvailableUnits + 1
	user.ReturnedBooks = append(user.ReturnedBooks, book)
	user.TakenBooks = remove(user.TakenBooks, book)

	s.bookRepository.UpdateUnits(book)
	s.userRepository.UpdateReturnedBooks(user, user.ReturnedBooks)
	s.userRepository.UpdateTakenBooks(user, user.TakenBooks)
}

func (s *userService) IsBookTakenByUser(email string, isbn string) bool {
	book := s.bookRepository.Find(isbn)
	if reflect.ValueOf(book).IsZero() {
		return false
	}

	user := s.userRepository.FindByEmail(email)
	if reflect.ValueOf(user).IsZero() {
		return false
	}

	return contains(user.TakenBooks, book)
}

func remove(slice []entities.Book, book entities.Book) []entities.Book {
	var s int
	for index, x := range slice {
		if x.Isbn == book.Isbn {
			s = index
			break
		}
	}
	return append(slice[:s], slice[s+1:]...)
}

func contains(slice []entities.Book, book entities.Book) bool {
	for _, x := range slice {
		if x.Isbn == book.Isbn {
			return true
		}
	}
	return false
}
