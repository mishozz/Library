package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/entities"
	"github.com/mishozz/Library/service"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type BookController interface {
	GetAll() []entities.Book
	GetByIsbn(isbn string) (entities.Book, error)
	Save(ctx *gin.Context) error
}

type bookController struct {
	service service.BookService
}

func NewController(service service.BookService) BookController {
	return &bookController{
		service: service,
	}
}

func (c *bookController) GetAll() []entities.Book {
	return c.service.FindAll()
}

func (c *bookController) Save(ctx *gin.Context) error {
	var book entities.Book
	err := ctx.ShouldBindJSON(&book)
	if err != nil {
		return err
	}
	if c.service.BookExists(book.Isbn) {
		return ConflictError
	}

	c.service.Save(book)
	return nil

}

func (c *bookController) GetByIsbn(isbn string) (book entities.Book, err error) {
	if !c.service.BookExists(isbn) {
		return book, NotFoundError
	}
	return c.service.FindByIsbn(isbn), nil
}
