package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/config"
	"github.com/mishozz/Library/controller"
	"github.com/mishozz/Library/repositories"
	"github.com/mishozz/Library/service"
	"github.com/mishozz/Library/utils"
)

const (
	port string = "8080"
)

var (
	db                                         = config.NewDatabaseConfig()
	bookRepository repositories.BookRepository = repositories.NewBookRepository(db)
	bookService    service.BookService         = service.NewBookService(bookRepository)
	bookController controller.BookController   = controller.NewBookController(bookService)
)

func main() {
	defer utils.CloseDB(db.Connection)

	server := gin.New()

	controller.HandleRequests(server, bookController)
	server.Run(":" + port)
}
