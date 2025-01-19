package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

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
			wantBody:       jsonError("Unauthorized: Token not found")},
		{
			name:           "Invalid token",
			req:            requestWithInvalidToken,
			wantStatusCode: 401,
			wantBody:       jsonError("Unauthorized: Invalid token"),
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
	r.GET("/protected", AuthMiddleware(), protectedRoute)

	return r
}

func normalRoute(c *gin.Context) {
	c.String(200, "normal")
}

func protectedRoute(c *gin.Context) {
	c.String(200, "protected")
}

func jsonError(
	message string,
) string {
	res := gin.H{
		"error": message,
	}

	b, _ := json.Marshal(res)

	return string(b)
}
