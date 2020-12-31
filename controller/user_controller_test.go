package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) FindByEmail(email string) entities.User {
	args := m.Called(email)
	return args.Get(0).(entities.User)
}

func (m *mockUserService) FindAll() []entities.User {
	args := m.Called()
	return args.Get(0).([]entities.User)
}

func (m *mockUserService) UserExists(email string) bool {
	args := m.Called(email)
	return args.Get(0).(bool)
}

func (m *mockUserService) TakeBook(user entities.User, book entities.Book) {
	m.Called(user, book)
}

func (m *mockUserService) ReturnBook(user entities.User, book entities.Book) {
	m.Called(user, book)
}

func (m *mockUserService) IsBookTakenByUser(email string, isbn string) bool {
	args := m.Called(email, isbn)
	return args.Get(0).(bool)
}

func Test_NewUserController(t *testing.T) {
	userService := &mockUserService{}
	bookService := &mockBookService{}
	userController := NewUserController(userService, bookService)
	assert.NotNil(t, userController.userService)
	assert.NotNil(t, userController.bookService)
}

func Test_UserController_GetAll(t *testing.T) {
	book := entities.Book{
		Isbn:           "test",
		Author:         "test",
		Title:          "test",
		AvailableUnits: 3,
	}
	expectedUsers := []entities.User{{
		Email:         "email",
		TakenBooks:    []entities.Book{book},
		ReturnedBooks: []entities.Book{book},
	}}

	mockService := func(m *mockUserService) *mockUserService {
		m.On("FindAll").Return(expectedUsers)
		return m
	}
	mockBookService := &mockBookService{}
	mockUserService := &mockUserService{}

	userController := NewUserController(mockService(mockUserService), mockBookService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	userController.GetAll(c)

	var users []entities.User
	err := json.Unmarshal(w.Body.Bytes(), &users)
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, expectedUsers, users)
	assert.Equal(t, http.StatusOK, w.Code)

	mockUserService.AssertExpectations(t)
}

func Test_UserController_GetByEmail(t *testing.T) {
	tests := []struct {
		name        string
		mockService func(m *mockUserService) *mockUserService
		respBody    gin.H
		respStatus  int
	}{{
		name: "success",
		mockService: func(m *mockUserService) *mockUserService {
			m.On("UserExists", "email").Return(true)
			m.On("FindByEmail", "email").Return(entities.User{
				Email: "email",
			})
			return m
		},
		respBody:   gin.H{"Email": "email", "Returned_books": interface{}(nil), "Taken_books": interface{}(nil)},
		respStatus: 200,
	}, {
		name: "user does not exist",
		mockService: func(m *mockUserService) *mockUserService {
			m.On("UserExists", "email").Return(false)
			return m
		},
		respBody:   gin.H{ERROR_MESSAGE: USER_NOT_FOUND},
		respStatus: 404,
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockBookService := &mockBookService{}
			mockUserService := &mockUserService{}

			userController := NewUserController(tt.mockService(mockUserService), mockBookService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = append(c.Params, gin.Param{Key: "email", Value: "email"})
			userController.GetByEmail(c)

			var actualBody gin.H
			err := json.Unmarshal(w.Body.Bytes(), &actualBody)
			if err != nil {
				t.FailNow()
			}

			assert.Equal(t, tt.respBody, actualBody)
			assert.Equal(t, tt.respStatus, w.Code)
			mockUserService.AssertExpectations(t)
		})
	}
}

func Test_UserController_TakeBook(t *testing.T) {
	book := entities.Book{
		Isbn:           "test",
		Author:         "test",
		Title:          "test",
		AvailableUnits: 3,
	}
	user := entities.User{
		Email: "email",
	}
	tests := []struct {
		name            string
		mockBookService func(m *mockBookService) *mockBookService
		mockUserService func(m *mockUserService) *mockUserService
		respBody        gin.H
		respStatus      int
	}{{
		name: "success",
		mockBookService: func(m *mockBookService) *mockBookService {
			m.On("FindByIsbn", "test").Return(book)
			m.On("BookExists", "test").Return(true)
			return m
		},
		mockUserService: func(m *mockUserService) *mockUserService {
			m.On("UserExists", "email").Return(true)
			m.On("FindByEmail", "email").Return(user)
			m.On("TakeBook", user, book)
			return m
		},
		respBody:   gin.H{MESSAGE: "Book successfully taken"},
		respStatus: 201,
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockBook := &mockBookService{}
			mockUser := &mockUserService{}

			userController := NewUserController(tt.mockUserService(mockUser), tt.mockBookService(mockBook))

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = append(c.Params, gin.Param{Key: "email", Value: "email"}, gin.Param{Key: "isbn", Value: "test"})
			userController.TakeBook(c)

			var actualBody gin.H
			err := json.Unmarshal(w.Body.Bytes(), &actualBody)
			if err != nil {
				t.FailNow()
			}

			assert.Equal(t, tt.respBody, actualBody)
			assert.Equal(t, tt.respStatus, w.Code)
			mockBook.AssertExpectations(t)
			mockUser.AssertExpectations(t)
		})
	}
}
