package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/entities"
	"github.com/stretchr/testify/mock"
)

type mockBookController struct {
	mock.Mock
}

func (m *mockBookController) GetAll() []entities.Book {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]entities.Book)
}

func (m *mockBookController) Save(ctx *gin.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockBookController) GetByIsbn(isbn string) (book entities.Book, err error) {
	args := m.Called(isbn)
	if args.Get(0) == nil {
		return
	}
	return args.Get(0).(entities.Book), args.Error(1)
}

// func Test_getAllBooks(t *testing.T) {
// 	mockBookControllerClient := func(m *mockBookController) *mockBookController {
// 		m.On("GetAll").Return([]entities.Book{{
// 			Isbn:           "test",
// 			Author:         "test-author",
// 			Title:          "test-title",
// 			AvailableUnits: 1,
// 		}})
// 		return m
// 	}
// 	t.Run("success", func(t *testing.T) {
// 		m := &mockBookController{}
// 		ctx := gin.Context{}
// 		getAllBooks(&ctx, mockBookControllerClient(m))
// 		fmt.Printf(ctx.ContentType())
// 		m.AssertExpectations(t)
// 	})
// }
