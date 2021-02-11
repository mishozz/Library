package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/auth"
	"github.com/mishozz/Library/entities"
	"github.com/mishozz/Library/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockSignin struct{}

var signIn func(auth.AuthDetails) (string, error)

func (fs *mockSignin) SignIn(authD auth.AuthDetails) (string, error) {
	return signIn(authD)
}

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
	signIn = func(auth.AuthDetails) (string, error) {
		return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX3V1aWQiOiI4M2IwOTYxMi05ZGZjLTRjMWQtOGY3ZC1hNTg5YWNlYzcwODEiLCJhdXRob3JpemVkIjp0cnVlLCJ1c2VyX2lkIjo1fQ.otegNS-W9OE8RsqGtiyJRCB-H0YXBygNXP91qeCPdF8", nil
	}
	service.Authorize = &mockSignin{}

	tests := []struct {
		name            string
		mockUserService func(m *mockUserService) *mockUserService
		mockAuthRepo    func(m *mockAuthRepo) *mockAuthRepo
		input           entities.User
		token           string
		statusCode      int
	}{{
		name: "success",
		mockUserService: func(m *mockUserService) *mockUserService {
			m.On("FindByEmail", "email").Return(entities.User{
				Model: gorm.Model{
					ID: 1,
				},
				Email:    "email",
				Password: "$2a$14$lFwBpX3h15NHVhjwab9wSO3Crlf9sQWFFZI7DFpLJH8mH4av9dWH6",
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
			Email:    "email",
			Password: "123",
		},
		token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX3V1aWQiOiI4M2IwOTYxMi05ZGZjLTRjMWQtOGY3ZC1hNTg5YWNlYzcwODEiLCJhdXRob3JpemVkIjp0cnVlLCJ1c2VyX2lkIjo1fQ.otegNS-W9OE8RsqGtiyJRCB-H0YXBygNXP91qeCPdF8",
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

			r.POST("/login", loginController.Login)

			reqBody, err := json.Marshal(tt.input)
			if err != nil {
				t.FailNow()
			}

			c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))

			r.ServeHTTP(w, c.Request)

			var actualBody string
			err = json.Unmarshal(w.Body.Bytes(), &actualBody)
			if err != nil {
				t.FailNow()
			}

			if w.Code == 200 {
				assert.Equal(t, tt.token, actualBody)
			}

			assert.Equal(t, tt.statusCode, w.Code)
			mockAuth.AssertExpectations(t)
			mockUserService.AssertExpectations(t)
		})
	}
}
func Test_LoginController_Logout(t *testing.T) {
	tests := []struct {
		name            string
		mockUserService func(m *mockUserService) *mockUserService
		mockAuthRepo    func(m *mockAuthRepo) *mockAuthRepo
		requestToken    string
		statusCode      int
	}{{
		name: "success",
		mockAuthRepo: func(m *mockAuthRepo) *mockAuthRepo {
			m.On("DeleteAuth", mock.Anything).Return(nil)
			return m
		},
		mockUserService: func(m *mockUserService) *mockUserService {
			return m
		},
		requestToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX3V1aWQiOiJjMmUxYjBjMy00ZGRjLTQ0NjUtYWVkNC1iNGE2NDM5NzI4M2MiLCJhdXRob3JpemVkIjp0cnVlLCJ1c2VyX2lkIjoxfQ.FWbfdhEJeK7mjZ-lWvs9scuyUrSKPrC4xafUoEqkduc",
		statusCode:   200,
	}, {
		name: "unauthorized",
		mockAuthRepo: func(m *mockAuthRepo) *mockAuthRepo {
			return m
		},
		mockUserService: func(m *mockUserService) *mockUserService {
			return m
		},
		requestToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX3V1aWQiOiJjMmUxYjBjMy00ZGRjLTQ0NjUtYWVkNC1iNGE2NDM5NzI4M2MiLCJhdXRob3JpemVkIjp0cnVlLCJ1c2VyX2lkIjoxfQ.FWbfdhEJeK7mjZ-lWvs9scuyUrSKPrC4xafUoEqkducxx",
		statusCode:   401,
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := &mockAuthRepo{}
			mockUserService := &mockUserService{}
			loginController := NewLoginController(tt.mockAuthRepo(mockAuth), tt.mockUserService(mockUserService))

			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			r.POST("/logout", loginController.LogOut)

			c.Request, _ = http.NewRequest(http.MethodPost, "/logout", nil)
			tokenString := fmt.Sprintf("Bearer %v", tt.requestToken)
			c.Request.Header.Set("Authorization", tokenString)
			r.ServeHTTP(w, c.Request)

			var actualBody gin.H
			err := json.Unmarshal(w.Body.Bytes(), &actualBody)
			if err != nil {
				t.FailNow()
			}
			assert.Equal(t, tt.statusCode, w.Code)
			mockAuth.AssertExpectations(t)
			mockUserService.AssertExpectations(t)
		})
	}
}
