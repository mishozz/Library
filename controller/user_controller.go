package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/service"
	"github.com/mishozz/Library/utils"
)

const (
	USER_NOT_FOUND     = "User not found"
	NO_AVAILABLE_UNITS = "This book has not available copies"
	BOOK_ALREADY_TAKEN = "This book is already taken"
	BOOK_IS_NOT_TAKEN  = "This book is not taken"
	MESSAGE            = "message"
)

func HandleUserRequests(server *gin.Engine, userController UserController, loginController LoginController) {
	apiRoutes := server.Group(LIBRARY_API_V1)
	{
		apiRoutes.POST("login", func(c *gin.Context) {
			loginController.Login(c)
		})

		apiRoutes.POST("logout", func(c *gin.Context) {
			loginController.Login(c)
		})

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
	books, err := c.userService.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Internal error")
		return
	}
	ctx.JSON(http.StatusOK, books)
}

func (c *userController) GetByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	user, err := c.userService.FindByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{ERROR_MESSAGE: USER_NOT_FOUND})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *userController) TakeBook(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	email := ctx.Param("email")

	user, err := c.userService.FindByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{ERROR_MESSAGE: USER_NOT_FOUND})
		return
	}

	book, err := c.bookService.FindByIsbn(isbn)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{ERROR_MESSAGE: BOOK_NOT_FOUND})
		return
	}

	if book.AvailableUnits <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{ERROR_MESSAGE: NO_AVAILABLE_UNITS})
		return
	}
	if utils.Contains(user.TakenBooks, book) {
		ctx.JSON(http.StatusBadRequest, gin.H{ERROR_MESSAGE: BOOK_ALREADY_TAKEN})
		return
	}

	err = c.userService.TakeBook(user, book)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{MESSAGE: "unable to take book"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{MESSAGE: "Book successfully taken"})

}

func (c *userController) ReturnBook(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	email := ctx.Param("email")

	if !c.userService.IsBookTakenByUser(email, isbn) {
		ctx.JSON(http.StatusNotFound, gin.H{ERROR_MESSAGE: BOOK_IS_NOT_TAKEN})
		return
	}

	book, err := c.bookService.FindByIsbn(isbn)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{ERROR_MESSAGE: BOOK_NOT_FOUND})
		return
	}
	user, err := c.userService.FindByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{ERROR_MESSAGE: USER_NOT_FOUND})
		return
	}
	err = c.userService.ReturnBook(user, book)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{ERROR_MESSAGE: "unable to return book"})
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{MESSAGE: "Book successfuly returned"})

}
