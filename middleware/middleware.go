package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mishozz/Library/auth"
)

const (
	ADMIN = "Admin"
	USER  = "User"
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

func TokenRoleMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "You need to be authorized to access this route")
			c.Abort()
			return
		}
		tokenAuth, err := auth.ExtractTokenAuth(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		if tokenAuth.Role != role {
			c.JSON(http.StatusForbidden, fmt.Sprintf("This route is forbidenn for %s", tokenAuth.Role))
			c.Abort()
			return
		}
		c.Next()
	}
}
