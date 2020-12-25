package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/config"
	"github.com/mishozz/Library/controller"
	"github.com/mishozz/Library/repositories"
	"github.com/mishozz/Library/service"
)

const (
	port string = "8080"
)

var (
	db                                         = config.NewDatabaseConfig()
	bookRepository repositories.BookRepository = repositories.NewBookRepository(db)
	bookService    service.BookService         = service.NewBookService(bookRepository)
	bookController controller.BookController   = controller.NewController(bookService)
)

func main() {
	defer config.CloseDB(db)

	server := gin.New()

	controller.HandleRequests(server, bookController)
	server.Run(":" + port)
}
