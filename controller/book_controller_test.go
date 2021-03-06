package controller

import (
	"bytes"
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

const (
	PORT = "8080"
)

type mockBookService struct {
	mock.Mock
}

func (m *mockBookService) Save(book entities.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *mockBookService) FindAll() ([]entities.Book, error) {
	args := m.Called()
	return args.Get(0).([]entities.Book), args.Error(1)
}

func (m *mockBookService) FindByIsbn(isbn string) (entities.Book, error) {
	args := m.Called(isbn)
	return args.Get(0).(entities.Book), args.Error(1)
}

func (m *mockBookService) Delete(isbn string) error {
	args := m.Called(isbn)
	return args.Error(0)
}

func (m *mockBookService) IsBookTaken(isbn string) bool {
	args := m.Called(isbn)
	return args.Get(0).(bool)
}

func Test_NewBookController(t *testing.T) {
	service := &mockBookService{}
	bookController := NewBookController(service)
	assert.NotNil(t, bookController.service)
}

func Test_BookController_GetAll(t *testing.T) {
	expectedBooks := []entities.Book{
		{
			Isbn:           "test",
			Author:         "test",
			Title:          "test",
			AvailableUnits: 1,
		},
	}

	mockService := func(m *mockBookService) *mockBookService {
		m.On("FindAll").Return([]entities.Book{
			{
				Isbn:           "test",
				Author:         "test",
				Title:          "test",
				AvailableUnits: 1,
			},
		}, nil)
		return m
	}

	mock := &mockBookService{}
	controller := NewBookController(mockService(mock))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	controller.GetAll(c)

	var books []entities.Book
	err := json.Unmarshal(w.Body.Bytes(), &books)
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, expectedBooks, books)
	assert.Equal(t, http.StatusOK, w.Code)

	mock.AssertExpectations(t)
}

func Test_BookController_Save(t *testing.T) {
	validBook := entities.Book{
		Isbn:           "test",
		Author:         "test",
		Title:          "test",
		AvailableUnits: 1,
	}
	invalidBook := entities.Book{
		Author:         "test",
		Title:          "test",
		AvailableUnits: 1,
	}

	tests := []struct {
		name            string
		input           entities.Book
		mockBookService func(m *mockBookService) *mockBookService
		respBody        gin.H
		statusCode      int
	}{{
		name: "success",
		mockBookService: func(m *mockBookService) *mockBookService {
			m.On("Save", validBook).Return(nil)
			return m
		},
		input:      validBook,
		statusCode: http.StatusCreated,
		respBody:   gin.H{message: saveSuccess},
	}, {
		name: "invalid book",
		mockBookService: func(m *mockBookService) *mockBookService {
			return m
		},
		input:      invalidBook,
		statusCode: http.StatusUnprocessableEntity,
		respBody:   gin.H{errorMessage: invalidRequest},
	}, {
		name: "book already exists",
		mockBookService: func(m *mockBookService) *mockBookService {
			m.On("Save", validBook).Return(errors.New("book exists"))
			return m
		},
		input:      validBook,
		statusCode: http.StatusConflict,
		respBody:   gin.H{errorMessage: bookConflict},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockBookService{}
			controller := NewBookController(tt.mockBookService(mock))

			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			r.POST("/books", controller.Save)

			reqBody, err := json.Marshal(tt.input)
			if err != nil {
				t.FailNow()
			}

			c.Request, _ = http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(reqBody))

			r.ServeHTTP(w, c.Request)

			var actualBody gin.H
			err = json.Unmarshal(w.Body.Bytes(), &actualBody)
			if err != nil {
				t.FailNow()
			}

			assert.Equal(t, tt.respBody, actualBody)
			assert.Equal(t, tt.statusCode, w.Code)
			mock.AssertExpectations(t)
		})
	}

}

func Test_BookController_GetByIsbn(t *testing.T) {
	validBook := entities.Book{
		Isbn:           "test",
		Author:         "test",
		Title:          "test",
		AvailableUnits: 1,
	}
	tests := []struct {
		name            string
		isbn            string
		mockBookService func(m *mockBookService) *mockBookService
		respBody        gin.H
		statusCode      int
	}{{
		name: "success",
		mockBookService: func(m *mockBookService) *mockBookService {
			m.On("FindByIsbn", "test").Return(validBook, nil)
			return m
		},
		isbn:       "test",
		statusCode: http.StatusOK,
		respBody:   gin.H{"Author": "test", "AvailableUnits": float64(1), "Isbn": "test", "Title": "test"},
	}, {
		name: "invalid book",
		mockBookService: func(m *mockBookService) *mockBookService {
			m.On("FindByIsbn", "test").Return(entities.Book{}, errors.New("Not found"))
			return m
		},
		isbn:       "test",
		statusCode: http.StatusNotFound,
		respBody:   gin.H{errorMessage: bookNotFound},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockBookService{}
			controller := NewBookController(tt.mockBookService(mock))

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = append(c.Params, gin.Param{Key: "isbn", Value: "test"})
			controller.GetByIsbn(c)

			var actualBody gin.H
			err := json.Unmarshal(w.Body.Bytes(), &actualBody)
			if err != nil {
				t.FailNow()
			}

			assert.Equal(t, tt.respBody, actualBody)
			assert.Equal(t, tt.statusCode, w.Code)
			mock.AssertExpectations(t)
		})
	}
}

func Test_BookController_Delete(t *testing.T) {
	tests := []struct {
		name            string
		mockBookService func(m *mockBookService) *mockBookService
		statusCode      int
	}{{
		name: "success",
		mockBookService: func(m *mockBookService) *mockBookService {
			m.On("FindByIsbn", "test").Return(entities.Book{}, nil)
			m.On("IsBookTaken", "test").Return(false)
			m.On("Delete", "test").Return(nil)
			return m
		},
		statusCode: 204,
	}, {
		name: "book not found",
		mockBookService: func(m *mockBookService) *mockBookService {
			m.On("FindByIsbn", "test").Return(entities.Book{}, errors.New("not found"))
			return m
		},
		statusCode: 404,
	}, {
		name: "book is still take",
		mockBookService: func(m *mockBookService) *mockBookService {
			m.On("FindByIsbn", "test").Return(entities.Book{}, nil)
			m.On("IsBookTaken", "test").Return(true)
			return m
		},
		statusCode: 400,
	}, {
		name: "error while deleting",
		mockBookService: func(m *mockBookService) *mockBookService {
			m.On("FindByIsbn", "test").Return(entities.Book{}, nil)
			m.On("IsBookTaken", "test").Return(false)
			m.On("Delete", "test").Return(errors.New("internal error"))
			return m
		},
		statusCode: 500,
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockBookService{}
			controller := NewBookController(tt.mockBookService(mock))

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = append(c.Params, gin.Param{Key: "isbn", Value: "test"})
			controller.Delete(c)

			assert.Equal(t, tt.statusCode, w.Code)
			mock.AssertExpectations(t)
		})
	}
}
