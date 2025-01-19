package auth

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// defaultTestingSecret is the default secret used for testing
// DO NOT USE THIS IN PRODUCTION
// For production, use a secure secret from config
const defaultTestingSecret = "sampleSecret"

var (
	// ErrInvalidToken is returned when the token is invalid
	ErrInvalidToken = errors.New("Invalid token")

	// ErrMissingToken is returned when the token is missing
	ErrMissingToken = errors.New("Missing token")
)

// Move to utils
func ginError(error error) gin.H {
	return gin.H{
		"error": error.Error(),
	}
}

// JWTMiddleware is a middleware that checks for JWT token
// in the request header and validates it before allowing access
//
// If the token is invalid, it returns a 401 Unauthorized response
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, ginError(ErrMissingToken))
			c.Abort()
			return
		}

		// Extract token
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// TODO: Get secret from config
			return []byte(defaultTestingSecret), nil
		})

		if err != nil {
			c.JSON(401, ginError(err))

			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(401, ginError(ErrInvalidToken))
			c.Abort()
			return
		}

		c.Next()
	}
}
