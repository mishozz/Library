package repositories

import (
	"github.com/mishozz/Library/config"
	"github.com/mishozz/Library/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(user entities.User)
	FindByEmail(email string) entities.User
	FindAll() []entities.User
	UpdateTakenBooks(user entities.User, takenBooks []entities.Book)
	UpdateReturnedBooks(user entities.User, returnedBooks []entities.Book)
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

func (r *userRepository) FindByEmail(email string) entities.User {
	var user entities.User
	r.connection.Where("Email = ?", email).First(&user)
	r.connection.Model(&user).Association("TakenBooks").Find(&user.TakenBooks)
	r.connection.Model(&user).Association("ReturnedBooks").Find(&user.ReturnedBooks)
	return user
}

func (r *userRepository) FindAll() []entities.User {
	var users []entities.User
	r.connection.Preload("TakenBooks").Preload("ReturnedBooks").Find(&users)
	return users
}

func (r *userRepository) UpdateTakenBooks(user entities.User, takenBooks []entities.Book) {
	r.connection.Model(&user).Association("TakenBooks").Clear()
	r.connection.Model(&user).Association("TakenBooks").Append(takenBooks)
}

func (r *userRepository) UpdateReturnedBooks(user entities.User, returnedBooks []entities.Book) {
	r.connection.Model(&user).Association("ReturnedBooks").Clear()
	r.connection.Model(&user).Association("ReturnedBooks").Append(returnedBooks)
}
