package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/entities"
	"github.com/mishozz/Library/service"
	"github.com/pkg/errors"
)

const LIBRARY_API_V1 = "/library/api/v1/"

var (
	ConflictError = errors.New("Every book must have a unique ISBN!")
	NotFoundError = errors.New("Book not found")
)

func HandleRequests(server *gin.Engine, bookController BookController) {
	apiRoutes := server.Group(LIBRARY_API_V1)
	{
		apiRoutes.GET("/books", func(ctx *gin.Context) {
			bookController.GetAll(ctx)
		})

		apiRoutes.GET("/books/:isbn", func(ctx *gin.Context) {
			bookController.GetByIsbn(ctx)
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error_message": "Invalid request body"})
	}
	if c.service.BookExists(book.Isbn) {
		ctx.JSON(http.StatusConflict, gin.H{"error_message": ConflictError})
	}

	c.service.Save(book)
	ctx.JSON(http.StatusCreated, book)
}

func (c *bookController) GetByIsbn(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	if !c.service.BookExists(isbn) {
		ctx.JSON(http.StatusNotFound, gin.H{"error_message": NotFoundError})
	}
	ctx.JSON(200, c.service.FindByIsbn(isbn))
}
