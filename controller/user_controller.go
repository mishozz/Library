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
	BOOK_IS_NOT_TAKEN  = "This book is not taken"
	MESSAGE            = "message"
)

func HandleUserRequests(server *gin.Engine, userController UserController) {
	apiRoutes := server.Group(LIBRARY_API_V1)
	{
		apiRoutes.GET("users", func(ctx *gin.Context) {
			userController.GetAll(ctx)
		})
		apiRoutes.GET("users/:email", func(ctx *gin.Context) {
			userController.GetByEmail(ctx)
		})
		apiRoutes.POST("users/:email/:isbn", func(ctx *gin.Context) {
			userController.TakeBook(ctx)
		})
		apiRoutes.DELETE("users/:email/:isbn", func(ctx *gin.Context) {
			userController.ReturnBook(ctx)
		})
	}
}

type UserController interface {
	TakeBook(ctx *gin.Context)
	ReturnBook(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetByEmail(ctx *gin.Context)
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

func (c *userController) GetByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	if !c.userService.UserExists(email) {
		ctx.JSON(http.StatusNotFound, gin.H{ERROR_MESSAGE: USER_NOT_FOUND})
	} else {
		ctx.JSON(http.StatusOK, c.userService.FindByEmail(email))
	}
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
			ctx.JSON(http.StatusCreated, gin.H{MESSAGE: "Book successfully taken"})
		}
	}
}

func (c *userController) ReturnBook(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	email := ctx.Param("email")

	if !c.userService.IsBookTakenByUser(email, isbn) {
		ctx.JSON(http.StatusNotFound, gin.H{ERROR_MESSAGE: BOOK_IS_NOT_TAKEN})
	} else {
		book := c.bookService.FindByIsbn(isbn)
		user := c.userService.FindByEmail(email)

		c.userService.ReturnBook(user, book)
		ctx.JSON(http.StatusNoContent, gin.H{MESSAGE: "Book successfuly returned"})
	}
}
