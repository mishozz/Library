package utils

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CloseDB(connection *gorm.DB) error {
	zap.L().Debug("Closing db connection")

	sqlDB, err := connection.DB()
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
