package repositories

import (
	"fmt"
	"testing"

	"github.com/mishozz/Library/config"
	"github.com/mishozz/Library/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db config.Database

func clearDatabase() {
	deleteFromTables(db, "users", "books")
}

func deleteFromTables(db config.Database, tables ...string) {
	for _, table := range tables {
		db.Connection.Exec(fmt.Sprintf("DELETE FROM %s", table))
	}
}

func newTestDatabaseConnection() config.Database {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		errors.Wrap(err, "unable to open db connection")
	}
	db.AutoMigrate(&entities.Book{}, &entities.User{})

	return config.Database{
		Connection: db,
	}
}

func Test_NewBookRepository(t *testing.T) {
	db = config.Database{Connection: &gorm.DB{}}
	repo := NewBookRepository(db)
	assert.NotNil(t, repo.connection)
}

func Test_NewUserRepository(t *testing.T) {
	db = config.Database{Connection: &gorm.DB{}}
	repo := NewUserRepository(db)
	assert.NotNil(t, repo.connection)
}

func Test_BookRepository_Save_Find(t *testing.T) {
	db = newTestDatabaseConnection()
	defer clearDatabase()

	bookRepo := NewBookRepository(db)
	book := entities.Book{
		Isbn:           "test",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 3,
	}

	bookRepo.Save(book)

	found := bookRepo.Find(book.Isbn)

	assertEqualBooks(t, book, found)
}

func Test_BookRepository_FindAll(t *testing.T) {
	db = newTestDatabaseConnection()
	defer clearDatabase()

	bookRepo := NewBookRepository(db)
	book1 := entities.Book{
		Isbn:           "test1",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 3,
	}
	book2 := entities.Book{
		Isbn:           "test2",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 2,
	}

	bookRepo.Save(book1)
	bookRepo.Save(book2)

	books := bookRepo.FindAll()

	assertEqualBooks(t, book1, books[0])
	assertEqualBooks(t, book2, books[1])
}

func Test_UserRepository_Save_Find(t *testing.T) {
	db = newTestDatabaseConnection()
	defer clearDatabase()

	userRepo := NewUserRepository(db)
	takenBook := entities.Book{
		Isbn:           "test1",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 2,
	}
	returnedBook := entities.Book{
		Isbn:           "test2",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 2,
	}
	user := entities.User{
		Email: "email",
		TakenBooks: []entities.Book{
			takenBook,
		},
		ReturnedBooks: []entities.Book{
			returnedBook,
		},
	}
	userRepo.Save(user)

	found := userRepo.FindByEmail("email")
	assert.Equal(t, user.Email, found.Email)
	assertEqualBooks(t, user.TakenBooks[0], found.TakenBooks[0])
	assertEqualBooks(t, user.ReturnedBooks[0], found.ReturnedBooks[0])
}

func Test_UserRepository_FindAll(t *testing.T) {
	db = newTestDatabaseConnection()
	defer clearDatabase()

	takenBook := entities.Book{
		Isbn:           "test1",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 2,
	}
	returnedBook := entities.Book{
		Isbn:           "test2",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 2,
	}
	user1 := entities.User{
		Email:      "email1",
		TakenBooks: []entities.Book{takenBook},
		ReturnedBooks: []entities.Book{
			returnedBook,
		},
	}
	user2 := entities.User{
		Email: "email2",
		TakenBooks: []entities.Book{
			takenBook,
		},
		ReturnedBooks: []entities.Book{
			returnedBook,
		},
	}

	userRepo := NewUserRepository(db)

	userRepo.Save(user1)
	userRepo.Save(user2)

	users := userRepo.FindAll()

	assertEqualUsers(t, user1, users[0])
	assertEqualUsers(t, user2, users[1])
}

func Test_UserRepository_UpdateTakenBooks(t *testing.T) {
	db = newTestDatabaseConnection()
	defer clearDatabase()

	oldBook := entities.Book{
		Isbn:           "test1",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 2,
	}

	newBook := entities.Book{
		Isbn:           "test1",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 2,
	}

	user := entities.User{
		Email: "email",
		TakenBooks: []entities.Book{
			oldBook,
		},
	}

	userRepo := NewUserRepository(db)

	userRepo.Save(user)
	userRepo.UpdateTakenBooks(user, []entities.Book{newBook})

	user = userRepo.FindByEmail("email")
	assertEqualBooks(t, newBook, user.TakenBooks[0])
}

func Test_UserRepository_UpdateReturnedBooks(t *testing.T) {
	db = newTestDatabaseConnection()
	defer clearDatabase()

	oldBook := entities.Book{
		Isbn:           "test1",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 2,
	}

	newBook := entities.Book{
		Isbn:           "test1",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 2,
	}

	user := entities.User{
		Email: "email",
		ReturnedBooks: []entities.Book{
			oldBook,
		},
	}

	userRepo := NewUserRepository(db)

	userRepo.Save(user)
	userRepo.UpdateReturnedBooks(user, []entities.Book{newBook})

	user = userRepo.FindByEmail("email")
	assertEqualBooks(t, newBook, user.ReturnedBooks[0])
}

func assertEqualUsers(t *testing.T, expected entities.User, actual entities.User) {
	assert.Equal(t, expected.Email, actual.Email)
	for i, book := range actual.TakenBooks {
		assertEqualBooks(t, expected.TakenBooks[i], book)
	}
	for i, book := range actual.ReturnedBooks {
		assertEqualBooks(t, expected.ReturnedBooks[i], book)
	}
}

func assertEqualBooks(t *testing.T, expected entities.Book, actual entities.Book) {
	assert.Equal(t, expected.Isbn, actual.Isbn)
	assert.Equal(t, expected.Author, actual.Author)
	assert.Equal(t, expected.AvailableUnits, actual.AvailableUnits)
	assert.Equal(t, expected.Title, actual.Title)
}
