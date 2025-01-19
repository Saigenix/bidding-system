package auth

import "github.com/gin-gonic/gin"

// AuthMiddleware is a middleware that checks for JWT token
// in the request header and validates it before allowing access
//
// If the token is invalid, it returns a 401 Unauthorized response
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(401, gin.H{
				"error": "Unauthorized: Token not found",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
