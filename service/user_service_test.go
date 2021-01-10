package service

import (
	"errors"
	"testing"

	"github.com/mishozz/Library/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Save(user entities.User) {}

func (m *mockUserRepository) FindByEmail(email string) (entities.User, error) {
	args := m.Called(email)
	return args.Get(0).(entities.User), args.Error(1)
}
func (m *mockUserRepository) FindAll() ([]entities.User, error) {
	args := m.Called()
	return args.Get(0).([]entities.User), args.Error(1)
}

func (m *mockUserRepository) UpdateTakenBooks(user entities.User, takenBooks []entities.Book) error {
	args := m.Called(user, takenBooks)
	return args.Error(0)
}

func (m *mockUserRepository) UpdateReturnedBooks(user entities.User, returnedBooks []entities.Book) error {
	args := m.Called(user, returnedBooks)
	return args.Error(0)
}

func Test_NewUserService(t *testing.T) {
	userRepo := &mockUserRepository{}
	bookRepo := &mockBookRepository{}
	service := NewUserService(userRepo, bookRepo)
	assert.NotNil(t, service.bookRepository)
	assert.NotNil(t, service.userRepository)
}

func Test_UserService_FindByEmail(t *testing.T) {
	expectedUser := entities.User{
		Email: "test",
	}
	mockUserRepo := func(m *mockUserRepository) *mockUserRepository {
		m.On("FindByEmail", "test").Return(entities.User{
			Email: "test",
		}, nil)
		return m
	}
	mockUserRepository := &mockUserRepository{}
	mockBookRepository := &mockBookRepository{}
	service := NewUserService(mockUserRepo(mockUserRepository), mockBookRepository)
	user, _ := service.FindByEmail("test")
	assert.Equal(t, expectedUser, user)
	mockUserRepository.AssertExpectations(t)
}

func Test_UserService_FindAll(t *testing.T) {
	expectedUsers := []entities.User{entities.User{Email: "test1"}, entities.User{Email: "test2"}}
	mockUserRepo := func(m *mockUserRepository) *mockUserRepository {
		m.On("FindAll").Return([]entities.User{entities.User{Email: "test1"}, entities.User{Email: "test2"}}, nil)
		return m
	}
	mockUserRepository := &mockUserRepository{}
	mockBookRepository := &mockBookRepository{}
	service := NewUserService(mockUserRepo(mockUserRepository), mockBookRepository)
	users, _ := service.userRepository.FindAll()
	assert.Equal(t, expectedUsers, users)
	mockUserRepository.AssertExpectations(t)
}

func Test_UserService_TakeBook(t *testing.T) {
	book1 := entities.Book{
		Isbn:           "test",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 2,
	}

	book2 := entities.Book{
		Isbn:           "test",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 1,
	}
	tests := []struct {
		name         string
		user         entities.User
		mockUserRepo func(m *mockUserRepository) *mockUserRepository
		mockBookRepo func(m *mockBookRepository) *mockBookRepository
	}{{
		name: "book is not in returned books",
		mockUserRepo: func(m *mockUserRepository) *mockUserRepository {
			m.On("UpdateTakenBooks", entities.User{
				Email:      "email1",
				TakenBooks: []entities.Book{book2},
			}, []entities.Book{book2}).Return(nil).Once()
			return m
		},
		mockBookRepo: func(m *mockBookRepository) *mockBookRepository {
			m.On("UpdateUnits", book2).Return(nil).Once()
			return m
		},
		user: entities.User{Email: "email1"},
	}, {
		name: "book is in returned books",
		mockUserRepo: func(m *mockUserRepository) *mockUserRepository {
			m.On("UpdateTakenBooks", entities.User{
				Email:         "email1",
				TakenBooks:    []entities.Book{book2},
				ReturnedBooks: []entities.Book{book1},
			}, []entities.Book{book2}).Return(nil).Once()
			m.On("UpdateReturnedBooks", entities.User{
				Email:         "email1",
				TakenBooks:    []entities.Book{book2},
				ReturnedBooks: []entities.Book{},
			}, []entities.Book{}).Return(nil).Once()
			return m
		},
		mockBookRepo: func(m *mockBookRepository) *mockBookRepository {
			m.On("UpdateUnits", book2).Return(nil).Once()
			return m
		},
		user: entities.User{
			Email:         "email1",
			ReturnedBooks: []entities.Book{book1},
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := &mockUserRepository{}
			mockBookRepository := &mockBookRepository{}
			service := NewUserService(tt.mockUserRepo(mockUserRepository), tt.mockBookRepo(mockBookRepository))
			service.TakeBook(tt.user, book1)
			mockBookRepository.AssertExpectations(t)
			mockUserRepository.AssertExpectations(t)
		})
	}
}

func Test_UserService_ReturnBook(t *testing.T) {
	book1 := entities.Book{
		Isbn:           "test",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 1,
	}

	book2 := entities.Book{
		Isbn:           "test",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 2,
	}
	mockUserRepo := func(m *mockUserRepository) *mockUserRepository {
		m.On("UpdateTakenBooks",
			entities.User{
				Email:         "email1",
				TakenBooks:    []entities.Book{},
				ReturnedBooks: []entities.Book{book2},
			},
			[]entities.Book{}).Return(nil).Once()
		m.On("UpdateReturnedBooks",
			entities.User{
				Email:         "email1",
				TakenBooks:    []entities.Book{},
				ReturnedBooks: []entities.Book{book2},
			},
			[]entities.Book{book2}).Return(nil).Once()
		return m
	}
	mockBookRepo := func(m *mockBookRepository) *mockBookRepository {
		m.On("UpdateUnits", book2).Return(nil).Once()
		return m
	}
	mockUserRepository := &mockUserRepository{}
	mockBookRepository := &mockBookRepository{}
	service := NewUserService(mockUserRepo(mockUserRepository), mockBookRepo(mockBookRepository))
	service.ReturnBook(entities.User{Email: "email1", TakenBooks: []entities.Book{book1}}, book1)
	mockBookRepository.AssertExpectations(t)
	mockUserRepository.AssertExpectations(t)
}

func Test_UserService_IsBookTakenByUser(t *testing.T) {
	book := entities.Book{
		Isbn:           "test",
		Title:          "test",
		Author:         "test",
		AvailableUnits: 2,
	}

	tests := []struct {
		name         string
		mockUserRepo func(m *mockUserRepository) *mockUserRepository
		mockBookRepo func(m *mockBookRepository) *mockBookRepository
		expected     bool
	}{{
		name: "book is taken",
		mockUserRepo: func(m *mockUserRepository) *mockUserRepository {
			m.On("FindByEmail", "email").Return(entities.User{
				Email: "email",
				TakenBooks: []entities.Book{
					book,
				},
			}, nil)
			return m
		},
		mockBookRepo: func(m *mockBookRepository) *mockBookRepository {
			m.On("Find", "test").Return(book, nil)
			return m
		},
		expected: true,
	}, {
		name: "user does not exist",
		mockUserRepo: func(m *mockUserRepository) *mockUserRepository {
			m.On("FindByEmail", "email").Return(entities.User{}, nil)
			return m
		},
		mockBookRepo: func(m *mockBookRepository) *mockBookRepository {
			m.On("Find", "test").Return(book, nil)
			return m
		},
		expected: false,
	}, {
		name: "user has not taken this book",
		mockUserRepo: func(m *mockUserRepository) *mockUserRepository {
			m.On("FindByEmail", "email").Return(entities.User{
				Email: "email",
			}, nil)
			return m
		},
		mockBookRepo: func(m *mockBookRepository) *mockBookRepository {
			m.On("Find", "test").Return(book, nil)
			return m
		},
		expected: false,
	}, {
		name: "book does not exist",
		mockUserRepo: func(m *mockUserRepository) *mockUserRepository {
			return m
		},
		mockBookRepo: func(m *mockBookRepository) *mockBookRepository {
			m.On("Find", "test").Return(entities.Book{}, errors.New("not found"))
			return m
		},
		expected: false,
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := &mockUserRepository{}
			mockBookRepository := &mockBookRepository{}
			service := NewUserService(tt.mockUserRepo(mockUserRepository), tt.mockBookRepo(mockBookRepository))
			flag := service.IsBookTakenByUser("email", "test")
			assert.Equal(t, tt.expected, flag)
			mockBookRepository.AssertExpectations(t)
			mockUserRepository.AssertExpectations(t)
		})
	}
}
