package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/service"
	"github.com/stretchr/stew/slice"
)

const (
	USER_NOT_FOUND     = "User not found"
	NO_AVAILABLE_UNITS = "This book has not available copies"
	BOOK_ALREADY_TAKEN = "This book is already taken"
)

func HandleUserRequests(server *gin.Engine, userController UserController) {
	apiRoutes := server.Group(LIBRARY_API_V1)
	{
		apiRoutes.GET("users", func(ctx *gin.Context) {
			userController.GetAll(ctx)
		})
		apiRoutes.POST("users/:email/:isbn", func(ctx *gin.Context) {
			userController.TakeBook(ctx)
		})
	}
}

type UserController interface {
	TakeBook(ctx *gin.Context)
	GetAll(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
	bookService service.BookService
}

func NewUserController(userService service.UserService, bookService service.BookService) *userController {
	return &userController{
		userService: userService,
		bookService: bookService,
	}
}

func (c *userController) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.userService.FindAll())
}

func (c *userController) TakeBook(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	email := ctx.Param("email")

	if !c.userService.UserExists(email) {
		ctx.JSON(http.StatusNotFound, gin.H{ERROR_MESSAGE: USER_NOT_FOUND})
	} else if !c.bookService.BookExists(isbn) {
		ctx.JSON(http.StatusNotFound, gin.H{ERROR_MESSAGE: BOOK_NOT_FOUND})
	} else {
		book := c.bookService.FindByIsbn(isbn)
		user := c.userService.FindByEmail(email)

		if book.AvailableUnits <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{ERROR_MESSAGE: NO_AVAILABLE_UNITS})
		} else if slice.Contains(user.TakenBooks, book) {
			ctx.JSON(http.StatusBadRequest, gin.H{ERROR_MESSAGE: BOOK_ALREADY_TAKEN})
		} else {
			c.userService.TakeBook(user, book)
			ctx.JSON(http.StatusCreated, gin.H{"message": "Book successfully taken"})
		}
	}
}
