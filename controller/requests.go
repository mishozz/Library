package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleRequests(server *gin.Engine, bookController BookController) {
	apiRoutes := server.Group("/library/api/v1/")
	{
		apiRoutes.GET("/books", func(ctx *gin.Context) {
			ctx.JSON(200, bookController.FindAll())
			//ctx.JSON(200, gin.H{"message": "suceess find all"})
		})
		apiRoutes.POST("/books", func(ctx *gin.Context) {
			err := bookController.Save(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Success!"})
			}
		})
	}
}
