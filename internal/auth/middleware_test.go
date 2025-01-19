package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/golang-jwt/jwt/v5"
)

// sampleSecret used in testing for JWT
const sampleSecret = "sampleSecret"

// TestPingRoute is a sample testcase
func TestPingRoute(
	t *testing.T,
) {

	router := setupRouter()

	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/ping", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	assert.Equal(t, "pong", w.Body.String())

}

// TestNormalRoute tests the normal route
func TestNormalRoute(
	t *testing.T,
) {
	router := setupRouter()

	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/normal", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	assert.Equal(t, "normal", w.Body.String())
}

// TestProtectedRoute tests the protected route
func TestProtectedRoute(
	t *testing.T,
) {

	router := setupRouter()
	requestWithInvalidToken := httptest.NewRequest("GET", "/protected", nil)
	requestWithInvalidToken.Header.Add("Authorization", "Bearer someInvalidToken")

	requestWithValidToken := httptest.NewRequest("GET", "/protected", nil)
	tokUnsigned := jwt.New(jwt.SigningMethodHS256)
	token, err := tokUnsigned.SignedString([]byte(sampleSecret))
	if err != nil {
		t.Fatal(err)
	}
	requestWithValidToken.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	tests := []struct {
		name           string
		wantStatusCode int
		wantBody       interface{}
		req            *http.Request
	}{
		{
			name:           "No token",
			req:            httptest.NewRequest("GET", "/protected", nil),
			wantStatusCode: 401,
			wantBody:       jsonError(ErrMissingToken)},
		{
			name:           "Invalid token",
			req:            requestWithInvalidToken,
			wantStatusCode: 401,
			wantBody:       jsonError(fmt.Errorf("%s: token contains an invalid number of segments", jwt.ErrTokenMalformed)),
		},
		{
			name:           "Valid token",
			req:            requestWithValidToken,
			wantStatusCode: 200,
			wantBody:       "protected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, tt.req)

			assert.Equal(t, tt.wantStatusCode, w.Code)

			if tt.wantBody != "" {
				assert.Equal(t, tt.wantBody, w.Body.String())
			}
		})
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/normal", normalRoute)
	r.GET("/protected", JWTMiddleware(), protectedRoute)

	return r
}

func normalRoute(c *gin.Context) {
	c.String(200, "normal")
}

func protectedRoute(c *gin.Context) {
	c.String(200, "protected")
}

func jsonError(
	err error,
) string {
	res := gin.H{
		"error": err.Error(),
	}

	b, _ := json.Marshal(res)

	return string(b)
}
