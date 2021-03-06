package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/controller"
	"github.com/mishozz/Library/middleware"
)

const (
	libraryApiV1 = "/library/api/v1/"
	ADMIN        = "Admin"
	USER         = "User"
)

// HandleRequests handles all incoming http requests
func HandleRequests(server *gin.Engine, bookController controller.BookController, userController controller.UserController, loginController controller.LoginController) {
	apiRoutes := server.Group(libraryApiV1)
	{
		apiRoutes.GET("/books", middleware.TokenAuthMiddleware(), func(ctx *gin.Context) {
			bookController.GetAll(ctx)
		})

		apiRoutes.GET("/books/:isbn", middleware.TokenAuthMiddleware(), func(ctx *gin.Context) {
			bookController.GetByIsbn(ctx)
		})

		apiRoutes.DELETE("/books/:isbn", middleware.TokenRoleMiddleware(ADMIN), func(ctx *gin.Context) {
			bookController.Delete(ctx)
		})

		apiRoutes.POST("/books", middleware.TokenRoleMiddleware(ADMIN), func(ctx *gin.Context) {
			bookController.Save(ctx)
		})

		apiRoutes.POST("register", func(c *gin.Context) {
			loginController.Register(c)
		})

		apiRoutes.POST("login", func(c *gin.Context) {
			loginController.Login(c)
		})

		apiRoutes.POST("logout", func(c *gin.Context) {
			loginController.LogOut(c)
		})

		apiRoutes.GET("users", middleware.TokenAuthMiddleware(), func(ctx *gin.Context) {
			userController.GetAll(ctx)
		})
		apiRoutes.GET("users/:email", middleware.TokenAuthMiddleware(), func(ctx *gin.Context) {
			userController.GetByEmail(ctx)
		})
		apiRoutes.POST("users/:email/:isbn", middleware.TokenRoleMiddleware(USER), func(ctx *gin.Context) {
			userController.TakeBook(ctx)
		})
		apiRoutes.DELETE("users/:email/:isbn", middleware.TokenRoleMiddleware(USER), func(ctx *gin.Context) {
			userController.ReturnBook(ctx)
		})
	}
}
