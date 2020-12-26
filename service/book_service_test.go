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

func (m *mockBookRepository) Save(entities.Book) {}

func (m *mockBookRepository) Delete(book entities.Book) {}

func (m *mockBookRepository) FindAll() []entities.Book {
	args := m.Called()
	return args.Get(0).([]entities.Book)
}

func (m *mockBookRepository) Find(isbn string) entities.Book {
	args := m.Called(isbn)
	return args.Get(0).(entities.Book)
}

func Test_NewBookService(t *testing.T) {
	repo := &mockBookRepository{}
	service := NewBookService(repo)
	assert.NotNil(t, service.repository)
}

func Test_BookService_Save(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := bookService{
			repository: &mockBookRepository{},
		}
		book := entities.Book{}
		service.Save(book)
	})
}

func Test_BookService_FindAll(t *testing.T) {
	tests := []struct {
		name         string
		mockBookRepo func(m *mockBookRepository) *mockBookRepository
		expected     []entities.Book
	}{{
		name: "find empty slice",
		mockBookRepo: func(m *mockBookRepository) *mockBookRepository {
			m.On("FindAll").Return([]entities.Book{})
			return m
		},
		expected: []entities.Book{},
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
			})
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
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := &mockBookRepository{}
			service := bookService{
				repository: tt.mockBookRepo(m),
			}
			books := service.FindAll()
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
		})
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
	book := service.FindByIsbn("test")
	assert.Equal(t, expectedBook, book)
	m.AssertExpectations(t)
}

func Test_BookService_BookExists(t *testing.T) {
	tests := []struct {
		name         string
		mockBookRepo func(m *mockBookRepository) *mockBookRepository
		expected     bool
	}{{
		name: "book exists",
		mockBookRepo: func(m *mockBookRepository) *mockBookRepository {
			m.On("Find", "test").Return(entities.Book{
				Isbn:           "test",
				Author:         "test",
				Title:          "test",
				AvailableUnits: 1,
			})
			return m
		},
		expected: true,
	}, {
		name: "book doesn't exist",
		mockBookRepo: func(m *mockBookRepository) *mockBookRepository {
			m.On("Find", "test").Return(entities.Book{})
			return m
		},
		expected: false,
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := &mockBookRepository{}
			service := &bookService{
				repository: tt.mockBookRepo(m),
			}
			flag := service.BookExists("test")
			assert.Equal(t, tt.expected, flag)
			m.AssertExpectations(t)
		})
	}
}
