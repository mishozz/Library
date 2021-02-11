package service

import (
	"github.com/mishozz/Library/entities"
	"github.com/mishozz/Library/repositories"
	"github.com/mishozz/Library/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindByEmail(email string) (entities.User, error)
	FindAll() ([]entities.User, error)
	TakeBook(user entities.User, book entities.Book) error
	ReturnBook(user entities.User, book entities.Book) error
	IsBookTakenByUser(email string, isbn string) bool
	Register(user entities.User) error
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

func (s *userService) FindByEmail(email string) (entities.User, error) {
	return s.userRepository.FindByEmail(email)
}

func (s *userService) FindAll() ([]entities.User, error) {
	return s.userRepository.FindAll()
}

func (s *userService) TakeBook(user entities.User, book entities.Book) error {
	book.AvailableUnits = book.AvailableUnits - 1
	user.TakenBooks = append(user.TakenBooks, book)

	err := s.bookRepository.UpdateUnits(book)
	if err != nil {
		return err
	}
	err = s.userRepository.UpdateTakenBooks(user, user.TakenBooks)
	if err != nil {
		return err
	}

	if utils.Contains(user.ReturnedBooks, book) {
		user.ReturnedBooks = utils.Remove(user.ReturnedBooks, book)
		err = s.userRepository.UpdateReturnedBooks(user, user.ReturnedBooks)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *userService) ReturnBook(user entities.User, book entities.Book) error {
	book.AvailableUnits = book.AvailableUnits + 1
	user.ReturnedBooks = append(user.ReturnedBooks, book)
	user.TakenBooks = utils.Remove(user.TakenBooks, book)

	err := s.bookRepository.UpdateUnits(book)
	if err != nil {
		return err
	}
	err = s.userRepository.UpdateReturnedBooks(user, user.ReturnedBooks)
	if err != nil {
		return err
	}
	err = s.userRepository.UpdateTakenBooks(user, user.TakenBooks)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) IsBookTakenByUser(email string, isbn string) bool {
	book, err := s.bookRepository.Find(isbn)
	if err != nil {
		return false
	}

	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return false
	}

	return utils.Contains(user.TakenBooks, book)
}

func (s *userService) Register(user entities.User) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	user.Role = "User"
	return s.userRepository.Save(user)
}
