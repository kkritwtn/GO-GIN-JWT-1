package middleware

import (
	"fmt"
	"go-jwt-crud/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		claims, err := utils.ValidateJWT(auth)
		if err != nil || claims.Role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized admin"})
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Next()
	}
}

func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		claims, err := utils.ValidateJWT(auth)
		fmt.Println("claims", claims)
		if err != nil || claims.Role != "user" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Set("user_id", claims.ID)
		c.Next()
	}
}
