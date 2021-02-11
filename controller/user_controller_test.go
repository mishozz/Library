package controller

import (
	"encoding/json"
	"errors"
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

func (m *mockUserService) FindByEmail(email string) (entities.User, error) {
	args := m.Called(email)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *mockUserService) FindAll() ([]entities.User, error) {
	args := m.Called()
	return args.Get(0).([]entities.User), args.Error(1)
}

func (m *mockUserService) TakeBook(user entities.User, book entities.Book) error {
	args := m.Called(user, book)
	return args.Error(0)
}

func (m *mockUserService) ReturnBook(user entities.User, book entities.Book) error {
	args := m.Called(user, book)
	return args.Error(0)
}

func (m *mockUserService) IsBookTakenByUser(email string, isbn string) bool {
	args := m.Called(email, isbn)
	return args.Get(0).(bool)
}

func (m *mockUserService) Register(user entities.User) error {
	args := m.Called(user)
	return args.Error(0)
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
		m.On("FindAll").Return(expectedUsers, nil)
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
			m.On("FindByEmail", "email").Return(entities.User{
				Email: "email",
			}, nil)
			return m
		},
		respBody:   gin.H{"Email": "email", "Returned_books": interface{}(nil), "Taken_books": interface{}(nil)},
		respStatus: 200,
	}, {
		name: "user does not exist",
		mockService: func(m *mockUserService) *mockUserService {
			m.On("FindByEmail", "email").Return(entities.User{}, errors.New("Not found"))
			return m
		},
		respBody:   gin.H{errorMessage: userNotFound},
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
	invalidBook := entities.Book{
		Isbn:           "test",
		Author:         "test",
		Title:          "test",
		AvailableUnits: 0,
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
			m.On("FindByIsbn", "test").Return(book, nil)
			return m
		},
		mockUserService: func(m *mockUserService) *mockUserService {
			m.On("FindByEmail", "email").Return(entities.User{
				Email: "email",
			}, nil)
			m.On("TakeBook", entities.User{
				Email: "email",
			}, book).Return(nil)
			return m
		},
		respBody:   gin.H{message: "Book successfully taken"},
		respStatus: 201,
	}, {
		name: "no available units",
		mockBookService: func(m *mockBookService) *mockBookService {
			m.On("FindByIsbn", "test").Return(invalidBook, nil)
			return m
		},
		mockUserService: func(m *mockUserService) *mockUserService {
			m.On("FindByEmail", "email").Return(entities.User{
				Email: "email",
			}, nil)
			return m
		},
		respBody:   gin.H{errorMessage: noAvailableUnits},
		respStatus: 400,
	}, {
		name: "book is already taken by this user",
		mockBookService: func(m *mockBookService) *mockBookService {
			m.On("FindByIsbn", "test").Return(book, nil)
			return m
		},
		mockUserService: func(m *mockUserService) *mockUserService {
			m.On("FindByEmail", "email").Return(entities.User{
				Email:      "email",
				TakenBooks: []entities.Book{book},
			}, nil)
			return m
		},
		respBody:   gin.H{errorMessage: bookAlreadyTaken},
		respStatus: 400,
	}, {
		name: "user doesnt exists",
		mockUserService: func(m *mockUserService) *mockUserService {
			m.On("FindByEmail", "email").Return(entities.User{}, errors.New("Not found"))
			return m
		},
		mockBookService: func(m *mockBookService) *mockBookService {
			return m
		},
		respBody:   gin.H{errorMessage: userNotFound},
		respStatus: 404,
	}, {
		name: "book doesnt exists",
		mockUserService: func(m *mockUserService) *mockUserService {
			m.On("FindByEmail", "email").Return(entities.User{}, nil)
			return m
		},
		mockBookService: func(m *mockBookService) *mockBookService {
			m.On("FindByIsbn", "test").Return(book, errors.New("Not found"))
			return m
		},
		respBody:   gin.H{errorMessage: bookNotFound},
		respStatus: 404,
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

func Test_UserController_ReturnBook(t *testing.T) {
	book := entities.Book{
		Isbn:           "test",
		Author:         "test",
		Title:          "test",
		AvailableUnits: 3,
	}
	tests := []struct {
		name            string
		mockBookService func(m *mockBookService) *mockBookService
		mockUserService func(m *mockUserService) *mockUserService
		respStatus      int
	}{{
		name: "success",
		mockBookService: func(m *mockBookService) *mockBookService {
			m.On("FindByIsbn", "test").Return(book, nil)
			return m
		},
		mockUserService: func(m *mockUserService) *mockUserService {
			m.On("IsBookTakenByUser", "email", "test").Return(true)
			m.On("FindByEmail", "email").Return(entities.User{
				Email:      "email",
				TakenBooks: []entities.Book{book},
			}, nil)
			m.On("ReturnBook", entities.User{
				Email:      "email",
				TakenBooks: []entities.Book{book},
			}, book).Return(nil)
			return m
		},
		respStatus: 204,
	}, {
		name: "book is not taken",
		mockBookService: func(m *mockBookService) *mockBookService {
			return m
		},
		mockUserService: func(m *mockUserService) *mockUserService {
			m.On("IsBookTakenByUser", "email", "test").Return(false)
			return m
		},
		respStatus: 404,
	}, {
		name: "user not found",
		mockUserService: func(m *mockUserService) *mockUserService {
			m.On("IsBookTakenByUser", "email", "test").Return(true)
			m.On("FindByEmail", "email").Return(entities.User{}, errors.New("Not found"))
			return m
		},
		mockBookService: func(m *mockBookService) *mockBookService {
			m.On("FindByIsbn", "test").Return(book, nil)
			return m
		},
		respStatus: 404,
	}, {
		name: "book not found",
		mockUserService: func(m *mockUserService) *mockUserService {
			m.On("IsBookTakenByUser", "email", "test").Return(true)
			return m
		},
		mockBookService: func(m *mockBookService) *mockBookService {
			m.On("FindByIsbn", "test").Return(entities.Book{}, errors.New("Not found"))
			return m
		},
		respStatus: 404,
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
			userController.ReturnBook(c)

			assert.Equal(t, tt.respStatus, w.Code)
			mockBook.AssertExpectations(t)
			mockUser.AssertExpectations(t)
		})
	}
}
