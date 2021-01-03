package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/auth"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "You need to be authorized to access this route")
			c.Abort()
			return
		}
		c.Next()
	}
}
