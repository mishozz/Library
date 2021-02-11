package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/auth"
	"github.com/mishozz/Library/entities"
	"github.com/mishozz/Library/repositories"
	"github.com/mishozz/Library/service"
	"golang.org/x/crypto/bcrypt"

	"net/http"
)

// LoginController is an interface with all the methods we need for the login controller
type LoginController interface {
	Login(c *gin.Context)
	LogOut(c *gin.Context)
	Register(c *gin.Context)
}

type loginController struct {
	authRepository repositories.AuthRepository
	userService    service.UserService
}

var (
	successullyRegister string = "registered successully"
	userConflict        string = "this user already exists"
	wrongPassword       string = "wrong password"
)

// NewLoginController creates a new instance of the login controller
func NewLoginController(authRepo repositories.AuthRepository, userService service.UserService) *loginController {
	return &loginController{
		authRepository: authRepo,
		userService:    userService,
	}
}

func (lc *loginController) Login(c *gin.Context) {
	var u entities.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	//check if the user exist:
	user, err := lc.userService.FindByEmail(u.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, userNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, wrongPassword)
		return
	}
	//since after the user logged out, we destroyed that record in the database so that same jwt token can't be used twice. We need to create the token again
	authData, err := lc.authRepository.CreateAuth(uint64(user.ID), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var authD auth.AuthDetails
	authD.UserId = authData.UserID
	authD.AuthUuid = authData.AuthUUID
	authD.Role = authData.Role

	token, loginErr := service.Authorize.SignIn(authD)
	if loginErr != nil {
		c.JSON(http.StatusForbidden, gin.H{errorMessage: "Please try to login later"})
		return
	}
	c.JSON(http.StatusOK, token)
}

func (lc *loginController) LogOut(c *gin.Context) {
	au, err := auth.ExtractTokenAuth(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{errorMessage: "unauthorized"})
		return
	}
	delErr := lc.authRepository.DeleteAuth(au)
	if delErr != nil {
		log.Println(delErr)
		c.JSON(http.StatusUnauthorized, gin.H{errorMessage: "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{message: "Successfully logged out"})
}

func (lc *loginController) Register(c *gin.Context) {
	var user entities.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{errorMessage: invalidRequest})
		return
	}
	err = lc.userService.Register(user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{errorMessage: userConflict})
		return
	}
	c.JSON(http.StatusCreated, gin.H{message: successullyRegister})
}
