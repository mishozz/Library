package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/auth"
	"github.com/mishozz/Library/entities"
	"github.com/mishozz/Library/repositories"
	"github.com/mishozz/Library/service"

	"net/http"
)

type LoginController interface {
	Login(c *gin.Context)
	LogOut(c *gin.Context)
}

type loginController struct {
	authRepository repositories.AuthRepository
	userService    service.UserService
}

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
	user := lc.userService.FindByEmail(u.Email)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, err.Error())
	// 	return
	// }
	//since after the user logged out, we destroyed that record in the database so that same jwt token can't be used twice. We need to create the token again
	authData, err := lc.authRepository.CreateAuth(uint64(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var authD auth.AuthDetails
	authD.UserId = authData.UserID
	authD.AuthUuid = authData.AuthUUID

	token, loginErr := service.Authorize.SignIn(authD)
	if loginErr != nil {
		c.JSON(http.StatusForbidden, "Please try to login later")
		return
	}
	c.JSON(http.StatusOK, token)
}

func (lc *loginController) LogOut(c *gin.Context) {
	au, err := auth.ExtractTokenAuth(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	delErr := lc.authRepository.DeleteAuth(au)
	if delErr != nil {
		log.Println(delErr)
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}
