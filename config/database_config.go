package config

import (
	"github.com/mishozz/Library/entities"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
	db.AutoMigrate(&entities.Book{}, &entities.User{})

	return Database{
		Connection: db,
	}
}

func CloseDB(db Database) error {
	zap.L().Debug("Closing db connection")

	sqlDB, err := db.Connection.DB()
	if err != nil {
		return errors.Wrap(err, "unable to get the sqlDB")
	}
	err = sqlDB.Close()
	if err != nil {
		return errors.Wrap(err, "unable to close db connection")
	}

	zap.L().Debug("Successfully closed the db connection")
	return nil
}
