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
}
