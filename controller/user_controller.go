package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/service"
	"github.com/mishozz/Library/utils"
)

const (
	userNotFound     = "User not found"
	noAvailableUnits = "This book has not available copies"
	bookAlreadyTaken = "This book is already taken"
	bookIsNotTaken   = "This book is not taken"
	message          = "message"
)

// UserController is interface with all the methods we need for the user controller
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

// NewUserController creates new instance of the user controller
func NewUserController(userService service.UserService, bookService service.BookService) *userController {
	return &userController{
		userService: userService,
		bookService: bookService,
	}
}

func (c *userController) GetAll(ctx *gin.Context) {
	users, err := c.userService.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Internal error")
		return
	}
	for _, user := range users {
		user.Password = ""
	}
	ctx.JSON(http.StatusOK, users)
}

func (c *userController) GetByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	user, err := c.userService.FindByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{errorMessage: userNotFound})
		return
	}
	user.Password = ""
	ctx.JSON(http.StatusOK, user)
}

func (c *userController) TakeBook(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	email := ctx.Param("email")

	user, err := c.userService.FindByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{errorMessage: userNotFound})
		return
	}

	book, err := c.bookService.FindByIsbn(isbn)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{errorMessage: bookNotFound})
		return
	}

	if book.AvailableUnits <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{errorMessage: noAvailableUnits})
		return
	}
	if utils.Contains(user.TakenBooks, book) {
		ctx.JSON(http.StatusBadRequest, gin.H{errorMessage: bookAlreadyTaken})
		return
	}

	err = c.userService.TakeBook(user, book)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{message: "unable to take book"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{message: "Book successfully taken"})

}

func (c *userController) ReturnBook(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	email := ctx.Param("email")

	if !c.userService.IsBookTakenByUser(email, isbn) {
		ctx.JSON(http.StatusNotFound, gin.H{errorMessage: bookIsNotTaken})
		return
	}

	book, err := c.bookService.FindByIsbn(isbn)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{errorMessage: bookNotFound})
		return
	}
	user, err := c.userService.FindByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{errorMessage: userNotFound})
		return
	}
	err = c.userService.ReturnBook(user, book)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{errorMessage: "unable to return book"})
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{message: "Book successfuly returned"})

}
