package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/entities"
	"github.com/mishozz/Library/service"
)

const (
	LIBRARY_API_V1  = "/library/api/v1/"
	SAVE_SUCCESS    = "successfully saved"
	INVALID_REQUEST = "Invalid request body"
	ERROR_MESSAGE   = "error message"
	BOOK_CONFLICT   = "Every book must have a unique ISBN!"
	BOOK_NOT_FOUND  = "Book not found"
)

func HandleBookRequests(server *gin.Engine, bookController BookController) {
	apiRoutes := server.Group(LIBRARY_API_V1)
	{
		apiRoutes.GET("/books", func(ctx *gin.Context) {
			bookController.GetAll(ctx)
		})

		apiRoutes.GET("/books/:isbn", func(ctx *gin.Context) {
			bookController.GetByIsbn(ctx)
		})

		apiRoutes.DELETE("/books/:isbn", func(ctx *gin.Context) {
			bookController.Delete(ctx)
		})

		apiRoutes.POST("/books", func(ctx *gin.Context) {
			bookController.Save(ctx)
		})
	}
}

type BookController interface {
	GetAll(ctx *gin.Context)
	GetByIsbn(ctx *gin.Context)
	Save(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type bookController struct {
	service service.BookService
}

func NewBookController(service service.BookService) *bookController {
	return &bookController{
		service: service,
	}
}

func (c *bookController) GetAll(ctx *gin.Context) {
	ctx.JSON(200, c.service.FindAll())
}

func (c *bookController) Save(ctx *gin.Context) {
	var book entities.Book
	err := ctx.ShouldBindJSON(&book)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ERROR_MESSAGE: INVALID_REQUEST})
	} else if c.service.BookExists(book.Isbn) {
		ctx.JSON(http.StatusConflict, gin.H{ERROR_MESSAGE: BOOK_CONFLICT})
	} else {
		c.service.Save(book)
		ctx.JSON(http.StatusCreated, gin.H{MESSAGE: SAVE_SUCCESS})
	}
}

func (c *bookController) GetByIsbn(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	if !c.service.BookExists(isbn) {
		ctx.JSON(http.StatusNotFound, gin.H{ERROR_MESSAGE: BOOK_NOT_FOUND})
	} else {
		ctx.JSON(200, c.service.FindByIsbn(isbn))
	}
}

func (c *bookController) Delete(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	if !c.service.BookExists(isbn) {
		ctx.JSON(http.StatusNotFound, gin.H{ERROR_MESSAGE: BOOK_NOT_FOUND})
	} else if c.service.IsBookTaken(isbn) {
		ctx.JSON(http.StatusBadRequest, gin.H{ERROR_MESSAGE: BOOK_ALREADY_TAKEN})
	} else {
		err := c.service.Delete(isbn)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{ERROR_MESSAGE: "error while deleting"})
		} else {
			ctx.JSON(204, gin.H{MESSAGE: "book deleted"})
		}
	}
}
