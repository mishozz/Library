package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/config"
	"github.com/mishozz/Library/controller"
	"github.com/mishozz/Library/repositories"
	"github.com/mishozz/Library/router"
	"github.com/mishozz/Library/service"
	"github.com/mishozz/Library/utils"
)

const (
	PORT string = "8080"
)

var (
	db                                         = config.NewDatabaseConfig()
	bookRepository repositories.BookRepository = repositories.NewBookRepository(db)
	userRepository repositories.UserRepository = repositories.NewUserRepository(db)
	authRepository repositories.AuthRepository = repositories.NewAuthRepository(db)

	bookService service.BookService = service.NewBookService(bookRepository)
	userService service.UserService = service.NewUserService(userRepository, bookRepository)

	bookController  controller.BookController  = controller.NewBookController(bookService)
	userController  controller.UserController  = controller.NewUserController(userService, bookService)
	loginController controller.LoginController = controller.NewLoginController(authRepository, userService)
)

func main() {
	defer utils.CloseDB(db.Connection)

	server := gin.New()

	router.HandleRequests(server, bookController, userController, loginController)

	server.Run(":" + PORT)
}
