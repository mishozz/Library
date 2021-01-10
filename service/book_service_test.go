package service

import (
	"testing"

	"github.com/mishozz/Library/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBookRepository struct {
	mock.Mock
}

func (m *mockBookRepository) Save(book entities.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *mockBookRepository) Delete(isbn string) error {
	args := m.Called(isbn)
	return args.Error(0)
}

func (m *mockBookRepository) IsBookTaken(isbn string) bool {
	args := m.Called(isbn)
	return args.Get(0).(bool)
}

func (m *mockBookRepository) FindAll() ([]entities.Book, error) {
	args := m.Called()
	return args.Get(0).([]entities.Book), args.Error(1)
}

func (m *mockBookRepository) UpdateUnits(book entities.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *mockBookRepository) Find(isbn string) (entities.Book, error) {
	args := m.Called(isbn)
	return args.Get(0).(entities.Book), args.Error(1)
}

func Test_NewBookService(t *testing.T) {
	repo := &mockBookRepository{}
	service := NewBookService(repo)
	assert.NotNil(t, service.repository)
}

func Test_BookService_Save(t *testing.T) {
	mockRepo := func(m *mockBookRepository) *mockBookRepository {
		m.On("Save", mock.Anything).Return(nil)
		return m
	}
	service := bookService{
		repository: mockRepo(&mockBookRepository{}),
	}
	book := entities.Book{}
	service.Save(book)
}

func Test_BookService_FindAll(t *testing.T) {
	tests := []struct {
		name         string
		mockBookRepo func(m *mockBookRepository) *mockBookRepository
		expected     []entities.Book
		err          error
	}{{
		name: "find empty slice",
		mockBookRepo: func(m *mockBookRepository) *mockBookRepository {
			m.On("FindAll").Return([]entities.Book{}, nil)
			return m
		},
		expected: []entities.Book{},
		err:      nil,
	}, {
		name: "find not empty slice",
		mockBookRepo: func(m *mockBookRepository) *mockBookRepository {
			m.On("FindAll").Return([]entities.Book{
				{
					Isbn:           "test",
					Author:         "test",
					Title:          "test",
					AvailableUnits: 1,
				},
			}, nil)
			return m
		},
		expected: []entities.Book{
			{
				Isbn:           "test",
				Author:         "test",
				Title:          "test",
				AvailableUnits: 1,
			},
		},
		err: nil,
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := &mockBookRepository{}
			service := bookService{
				repository: tt.mockBookRepo(m),
			}
			books, err := service.FindAll()
			if err != nil {
				assert.EqualError(t, err, tt.err.Error())
			}

			assert.Equal(t, tt.expected, books)
			m.AssertExpectations(t)
		})
	}
}

func Test_BookService_FindByIsbn(t *testing.T) {
	mockRepo := func(m *mockBookRepository) *mockBookRepository {
		m.On("Find", "test").Return(entities.Book{
			Isbn:           "test",
			Author:         "test",
			Title:          "test",
			AvailableUnits: 1,
		}, nil)
		return m
	}
	expectedBook := entities.Book{
		Isbn:           "test",
		Author:         "test",
		Title:          "test",
		AvailableUnits: 1,
	}
	m := &mockBookRepository{}
	service := bookService{
		repository: mockRepo(m),
	}
	book, _ := service.FindByIsbn("test")
	assert.Equal(t, expectedBook, book)
	m.AssertExpectations(t)
}

func Test_BookService_Delete(t *testing.T) {
	mockRepo := func(m *mockBookRepository) *mockBookRepository {
		m.On("Delete", "test").Return(nil)
		return m
	}

	m := &mockBookRepository{}
	service := bookService{
		repository: mockRepo(m),
	}
	err := service.Delete("test")
	assert.Nil(t, err)
	m.AssertExpectations(t)
}

func Test_BookService_IsBookTaken(t *testing.T) {
	mockRepo := func(m *mockBookRepository) *mockBookRepository {
		m.On("IsBookTaken", "test").Return(true)
		return m
	}

	m := &mockBookRepository{}
	service := bookService{
		repository: mockRepo(m),
	}
	err := service.IsBookTaken("test")
	assert.True(t, err)
	m.AssertExpectations(t)
}
