package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/entities"

	"github.com/mishozz/Library/repositories"
	"github.com/mishozz/Library/service"
)

const (
	saveSuccess    = "successfully saved"
	invalidRequest = "Invalid request body"
	errorMessage   = "error message"
	bookConflict   = "Every book must have a unique ISBN!"
	bookNotFound   = "Book not found"
	ADMIN          = "Admin"
	USER           = "User"
)

// BookController is an interface with all the methods we need for the book controller
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

// NewBookController creates a new instance of the book controller
func NewBookController(service service.BookService) *bookController {
	return &bookController{
		service: service,
	}
}

func (c *bookController) GetAll(ctx *gin.Context) {
	books, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(400, gin.H{errorMessage: "Internal error"})
		return
	}
	ctx.JSON(200, books)
}

func (c *bookController) Save(ctx *gin.Context) {
	var book entities.Book
	err := ctx.ShouldBindJSON(&book)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{errorMessage: invalidRequest})
		return
	}
	err = c.service.Save(book)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{errorMessage: bookConflict})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{message: saveSuccess})
}

func (c *bookController) GetByIsbn(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	book, err := c.service.FindByIsbn(isbn)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{errorMessage: bookNotFound})
	} else {
		ctx.JSON(200, book)
	}
}

func (c *bookController) Delete(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	_, err := c.service.FindByIsbn(isbn)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{errorMessage: bookNotFound})
		return
	}

	if c.service.IsBookTaken(isbn) {
		ctx.JSON(http.StatusBadRequest, gin.H{errorMessage: bookAlreadyTaken})
		return
	}

	err = c.service.Delete(isbn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{errorMessage: "error while deleting"})
	} else {
		ctx.JSON(204, gin.H{message: "book deleted"})
	}
}
