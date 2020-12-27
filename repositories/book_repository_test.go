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

func assertEqualBooks(t *testing.T, expected entities.Book, actual entities.Book) {
	assert.Equal(t, expected.Isbn, actual.Isbn)
	assert.Equal(t, expected.Author, actual.Author)
	assert.Equal(t, expected.AvailableUnits, actual.AvailableUnits)
	assert.Equal(t, expected.Title, actual.Title)
}

func Test_BookRepository_Save_Find(t *testing.T) {
	db := newTestDatabaseConnection()
	defer config.CloseDB(db)
	defer deleteFromTables(db, "users", "books")

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
	db := newTestDatabaseConnection()
	defer config.CloseDB(db)
	defer deleteFromTables(db, "users", "books")

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
