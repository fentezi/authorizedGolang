package middlewares

import (
	"github.com/gin-gonic/gin"
	"go-auth/utils"
)

func isAuthorized(c *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		claims, err := utils.ParserToken(cookie)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Set("role", claims.Role)
		c.Next()
	}
}
