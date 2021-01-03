package config

import (
	"github.com/mishozz/Library/entities"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Connection *gorm.DB
}

func NewDatabaseConfig() Database {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		errors.Wrap(err, "unable to open db connection")
	}
	db.AutoMigrate(&entities.Book{}, &entities.User{}, &entities.Auth{})

	return Database{
		Connection: db,
	}
}
