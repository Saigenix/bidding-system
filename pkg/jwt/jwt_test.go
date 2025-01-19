package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const sampleSecret = "sampleSecret"

// TestGenerateToken tests the GenerateToken function
func TestGenerateToken(t *testing.T) {

	jwtGenerator := NewJWTGenerator([]byte(sampleSecret))

	token, err := jwtGenerator.GenerateToken("Test User")

	assert.NoError(t, err, "Error generating token")

	assert.NotEmpty(t, token, "Token is empty")

}

// TestParseToken tests the ParseToken function
func TestParseToken(t *testing.T) {
	jwtGenerator := NewJWTGenerator([]byte(sampleSecret))

	token, err := jwtGenerator.GenerateToken("Test User")
	assert.NoError(t, err, "Error generating token")

	parsedToken, err := jwtGenerator.ParseToken(token)
	assert.NoError(t, err, "Error parsing token")

	assert.True(t, parsedToken.Valid, "Token is invalid")

	assert.Equal(t, "Test User", parsedToken.UserId())
}

// TestExpiredToken tests the ParseToken function with an expired token
func TestExpiredToken(t *testing.T) {
	jwtGenerator := NewJWTGenerator([]byte(sampleSecret))

	jwtGenerator.expiration = -1

	token, err := jwtGenerator.GenerateToken("Test User")
	assert.NoError(t, err, "Error generating token")

	parsedToken, err := jwtGenerator.ParseToken(token)
	assert.Error(t, err, "Error parsing token")

	assert.False(t, parsedToken.Valid, "Token is valid")
}

// TestCustomClaims tests the ParseToken function with custom claims
func TestCustomClaims(t *testing.T) {
	jwtGenerator := NewJWTGenerator([]byte(sampleSecret))

	token, err := jwtGenerator.GenerateToken(
		"Test User",
		WithClaims(map[string]interface{}{
			"role": "admin",
		}),
	)

	assert.NoError(t, err, "Error generating token")

	parsedToken, err := jwtGenerator.ParseToken(token)
	assert.NoError(t, err, "Error parsing token")

	assert.True(t, parsedToken.Valid, "Token is invalid")

	assert.Equal(t, "admin", parsedToken.claims["role"], "Role is not admin")
}

// ADDITIONAL TESTS
