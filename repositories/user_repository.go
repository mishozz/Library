package repositories

import (
	"github.com/mishozz/Library/config"
	"github.com/mishozz/Library/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(user entities.User)
	FindByEmail(email string) (entities.User, error)
	FindAll() ([]entities.User, error)
	UpdateTakenBooks(user entities.User, takenBooks []entities.Book) error
	UpdateReturnedBooks(user entities.User, returnedBooks []entities.Book) error
}

type userRepository struct {
	connection *gorm.DB
}

func NewUserRepository(db config.Database) *userRepository {
	return &userRepository{
		connection: db.Connection,
	}
}

func (r *userRepository) Save(user entities.User) {
	r.connection.Create(&user)
}

func (r *userRepository) FindByEmail(email string) (entities.User, error) {
	var user entities.User
	err := r.connection.Where("Email = ?", email).First(&user).Error
	if err != nil {
		return user, err
	}
	err = r.connection.Model(&user).Association("TakenBooks").Find(&user.TakenBooks)
	if err != nil {
		return user, err
	}
	err = r.connection.Model(&user).Association("ReturnedBooks").Find(&user.ReturnedBooks)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) FindAll() ([]entities.User, error) {
	var users []entities.User
	err := r.connection.Preload("TakenBooks").Preload("ReturnedBooks").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) UpdateTakenBooks(user entities.User, takenBooks []entities.Book) error {
	err := r.connection.Model(&user).Association("TakenBooks").Clear()
	if err != nil {
		return err
	}
	err = r.connection.Model(&user).Association("TakenBooks").Append(takenBooks)
	if err != nil {
		return err
	}
	return err
}

func (r *userRepository) UpdateReturnedBooks(user entities.User, returnedBooks []entities.Book) error {
	err := r.connection.Model(&user).Association("ReturnedBooks").Clear()
	if err != nil {
		return err
	}
	err = r.connection.Model(&user).Association("ReturnedBooks").Append(returnedBooks)
	if err != nil {
		return err
	}
	return nil
}
