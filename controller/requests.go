package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var ConflictErr = errors.New("Every book must have a unique ISBN!")

func HandleRequests(server *gin.Engine, bookController BookController) {
	apiRoutes := server.Group("/library/api/v1/")
	{
		apiRoutes.GET("/books", func(ctx *gin.Context) {
			ctx.JSON(200, bookController.GetAll())
			//ctx.JSON(200, gin.H{"message": "suceess find all"})
		})

		apiRoutes.GET("/books/:isbn", func(ctx *gin.Context) {
			isbn := ctx.Param("isbn")
			book, err := bookController.GetByIsbn(isbn)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("book with this isbn: %s does not exist", isbn)})
			}
			ctx.JSON(200, book)
		})

		apiRoutes.POST("/books", func(ctx *gin.Context) {
			err := bookController.Save(ctx)
			if err != nil {
				switch err {
				case ConflictErr:
					ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				default:
					ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				}
			} else {
				ctx.JSON(http.StatusCreated, gin.H{"message": "Success!"})
			}
		})
	}
}
