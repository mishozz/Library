package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/entities"

	"github.com/mishozz/Library/repositories"
	"github.com/mishozz/Library/service"
)

const (
	LIBRARY_API_V1  = "/library/api/v1/"
	SAVE_SUCCESS    = "successfully saved"
	INVALID_REQUEST = "Invalid request body"
	ERROR_MESSAGE   = "error message"
	BOOK_CONFLICT   = "Every book must have a unique ISBN!"
	BOOK_NOT_FOUND  = "Book not found"

	ADMIN = "Admin"
	USER  = "User"
)

type BookController interface {
	GetAll(ctx *gin.Context)
	GetByIsbn(ctx *gin.Context)
	Save(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type bookController struct {
	service        service.BookService
	authRepository repositories.AuthRepository
}

func NewBookController(service service.BookService) *bookController {
	return &bookController{
		service: service,
	}
}

func (c *bookController) GetAll(ctx *gin.Context) {
	books, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(400, gin.H{ERROR_MESSAGE: "Internal error"})
		return
	}
	ctx.JSON(200, books)
}

func (c *bookController) Save(ctx *gin.Context) {
	var book entities.Book
	err := ctx.ShouldBindJSON(&book)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{ERROR_MESSAGE: INVALID_REQUEST})
		return
	}
	err = c.service.Save(book)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{ERROR_MESSAGE: BOOK_CONFLICT})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{MESSAGE: SAVE_SUCCESS})
}

func (c *bookController) GetByIsbn(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	book, err := c.service.FindByIsbn(isbn)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{ERROR_MESSAGE: BOOK_NOT_FOUND})
	} else {
		ctx.JSON(200, book)
	}
}

func (c *bookController) Delete(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	_, err := c.service.FindByIsbn(isbn)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{ERROR_MESSAGE: BOOK_NOT_FOUND})
		return
	}

	if c.service.IsBookTaken(isbn) {
		ctx.JSON(http.StatusBadRequest, gin.H{ERROR_MESSAGE: BOOK_ALREADY_TAKEN})
		return
	}

	err = c.service.Delete(isbn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{ERROR_MESSAGE: "error while deleting"})
	} else {
		ctx.JSON(204, gin.H{MESSAGE: "book deleted"})
	}
}
