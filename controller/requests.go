package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	ConflictError = errors.New("Every book must have a unique ISBN!")
	NotFoundError = errors.New("Book not found")
)

func HandleRequests(server *gin.Engine, bookController BookController) {
	apiRoutes := server.Group("/library/api/v1/")
	{
		apiRoutes.GET("/books", func(ctx *gin.Context) {
			ctx.JSON(200, bookController.GetAll())
		})

		apiRoutes.GET("/books/:isbn", func(ctx *gin.Context) {
			isbn := ctx.Param("isbn")
			book, err := bookController.GetByIsbn(isbn)
			if err != nil {
				switch err {
				case NotFoundError:
					ctx.JSON(http.StatusNotFound, gin.H{"error_message": fmt.Sprintf("book with this isbn: %s does not exist", isbn)})
				default:
					ctx.JSON(http.StatusBadRequest, gin.H{"error_message": fmt.Sprintf("bad request")})
				}
			} else {
				ctx.JSON(200, book)
			}
		})

		apiRoutes.POST("/books", func(ctx *gin.Context) {
			err := bookController.Save(ctx)
			if err != nil {
				switch err {
				case ConflictError:
					ctx.JSON(http.StatusConflict, gin.H{"error_message": err.Error()})
				default:
					ctx.JSON(http.StatusBadRequest, gin.H{"error_message": "Invalid request body"})
				}
			} else {
				ctx.JSON(http.StatusCreated, gin.H{"message": "Success!"})
			}
		})
	}
}
