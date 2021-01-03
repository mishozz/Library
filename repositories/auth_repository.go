package repositories

import (
	"github.com/mishozz/Library/auth"
	"github.com/mishozz/Library/config"
	"github.com/mishozz/Library/entities"
	"github.com/twinj/uuid"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FetchAuth(*auth.AuthDetails) (*entities.Auth, error)
	DeleteAuth(*auth.AuthDetails) error
	CreateAuth(uint64) (*entities.Auth, error)
}

type authRepository struct {
	connection *gorm.DB
}

func NewAuthRepository(db config.Database) *authRepository {
	return &authRepository{
		connection: db.Connection,
	}
}

func (s *authRepository) FetchAuth(authD *auth.AuthDetails) (*entities.Auth, error) {
	au := &entities.Auth{}
	err := s.connection.Debug().Where("user_id = ? AND auth_uuid = ?", authD.UserId, authD.AuthUuid).Take(&au).Error
	if err != nil {
		return nil, err
	}
	return au, nil
}

//Once a user row in the auth table
func (s *authRepository) DeleteAuth(authD *auth.AuthDetails) error {
	au := &entities.Auth{}
	db := s.connection.Debug().Where("user_id = ? AND auth_uuid = ?", authD.UserId, authD.AuthUuid).Take(&au).Delete(&au)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

//Once the user signup/login, create a row in the auth table, with a new uuid
func (s *authRepository) CreateAuth(userId uint64) (*entities.Auth, error) {
	au := &entities.Auth{}
	au.AuthUUID = uuid.NewV4().String() //generate a new UUID each time
	au.UserID = userId
	err := s.connection.Debug().Create(&au).Error
	if err != nil {
		return nil, err
	}
	return au, nil
}
