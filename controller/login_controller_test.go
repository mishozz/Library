package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/auth"
	"github.com/mishozz/Library/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockAuthRepo struct {
	mock.Mock
}

func (m *mockAuthRepo) FetchAuth(auth *auth.AuthDetails) (*entities.Auth, error) {
	args := m.Called(auth)
	return args.Get(0).(*entities.Auth), args.Error(0)
}
func (m *mockAuthRepo) DeleteAuth(auth *auth.AuthDetails) error {
	args := m.Called(auth)
	return args.Error(0)
}
func (m *mockAuthRepo) CreateAuth(id uint64, role string) (auth *entities.Auth, err error) {
	args := m.Called(id, role)
	return args.Get(0).(*entities.Auth), args.Error(1)
}

func Test_LoginController_Login(t *testing.T) {
	tests := []struct {
		name            string
		mockUserService func(m *mockUserService) *mockUserService
		mockAuthRepo    func(m *mockAuthRepo) *mockAuthRepo
		input           entities.User
		statusCode      int
	}{{
		name: "success",
		mockUserService: func(m *mockUserService) *mockUserService {
			m.On("FindByEmail", "email").Return(entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Email: "email",
			}, nil)
			return m
		},
		mockAuthRepo: func(m *mockAuthRepo) *mockAuthRepo {
			m.On("CreateAuth", mock.Anything, mock.Anything).Return(&entities.Auth{
				ID:       1,
				UserID:   1,
				AuthUUID: "83b09612-9dfc-4c1d-8f7d-a589acec7081",
			}, nil)
			return m
		},
		input: entities.User{
			Email: "email",
		},
		statusCode: 200,
	}, {
		name: "user not found",
		mockUserService: func(m *mockUserService) *mockUserService {
			m.On("FindByEmail", "email").Return(entities.User{}, errors.New("not found"))
			return m
		},
		mockAuthRepo: func(m *mockAuthRepo) *mockAuthRepo {
			return m
		},
		input: entities.User{
			Email: "email",
		},
		statusCode: 404,
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := &mockAuthRepo{}
			mockUserService := &mockUserService{}
			loginController := NewLoginController(tt.mockAuthRepo(mockAuth), tt.mockUserService(mockUserService))

			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			r.POST("/books", loginController.Login)

			reqBody, err := json.Marshal(tt.input)
			if err != nil {
				t.FailNow()
			}

			c.Request, _ = http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(reqBody))

			r.ServeHTTP(w, c.Request)

			var actualBody string
			err = json.Unmarshal(w.Body.Bytes(), &actualBody)
			if err != nil {
				t.FailNow()
			}

			assert.Equal(t, tt.statusCode, w.Code)
			mockAuth.AssertExpectations(t)
			mockUserService.AssertExpectations(t)
		})
	}
}
