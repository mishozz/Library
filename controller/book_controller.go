package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/entities"
	"github.com/mishozz/Library/service"
)

type BookController interface {
	FindAll() []entities.Book
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

func (c *bookController) FindAll() []entities.Book {
	return c.service.FindAll()
}

func (c *bookController) Save(ctx *gin.Context) error {
	var book entities.Book
	err := ctx.ShouldBindJSON(&book)
	if err != nil {
		return err
	}
	// err = validate.Struct(video)
	// if err != nil {
	// 	return err
	// }
	c.service.Save(book)
	return nil

}
