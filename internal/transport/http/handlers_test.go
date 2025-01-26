package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
  "time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupTestRouter() *gin.Engine {
	router := gin.Default()
	server := Server{}
	RegisterHandlers(router, server)
	return router
}

func TestPostAuctions(t *testing.T) {
	router := SetupTestRouter()

	auction := AuctionCreate{
		Id:          "1",
		Title:       stringPtr("Test Auction"),
		Description: stringPtr("This is a test auction"),
	}

	auctionJSON, _ := json.Marshal(auction)

	req, _ := http.NewRequest("POST", "/auctions", bytes.NewBuffer(auctionJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	fmt.Println("Starting request")

	assert.Equal(t, http.StatusNotImplemented, w.Code)
	assert.Contains(t, w.Body.String(), "Not yet implemented")
}

func TestGetAuctions(t *testing.T) {
	router := SetupTestRouter()

	req, _ := http.NewRequest("GET", "/auctions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
	assert.Contains(t, w.Body.String(), "Not yet implemented")
}

func TestDeleteAuctionsAuctionId(t *testing.T) {
	router := SetupTestRouter()

	req, _ := http.NewRequest("DELETE", "/auctions/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
	assert.Contains(t, w.Body.String(), "Not yet implemented")
}

func TestGetAuctionsAuctionId(t *testing.T) {
	router := SetupTestRouter()

	req, _ := http.NewRequest("GET", "/auctions/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
	assert.Contains(t, w.Body.String(), "Not yet implemented")
}

func TestPutAuctionsAuctionId(t *testing.T) {
	router := SetupTestRouter()

	auction := AuctionCreate{
		Id:          "1",
		Title:       stringPtr("Updated Auction"),
		Description: stringPtr("This is an updated test auction"),
	}

	auctionJSON, _ := json.Marshal(auction)

	req, _ := http.NewRequest("PUT", "/auctions/1", bytes.NewBuffer(auctionJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
	assert.Contains(t, w.Body.String(), "Not yet implemented")
}

func TestPostBids(t *testing.T) {
	router := SetupTestRouter()

	bid := Bid{
		Amount:    100.0,
		AuctionId: "1",
		UserId:    "user1",
	}

	bidJSON, _ := json.Marshal(bid)

	req, _ := http.NewRequest("POST", "/bids", bytes.NewBuffer(bidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
	assert.Contains(t, w.Body.String(), "Not yet implemented")
}

func TestGetBidsBidId(t *testing.T) {
	router := SetupTestRouter()

	req, _ := http.NewRequest("GET", "/bids/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
	assert.Contains(t, w.Body.String(), "Not yet implemented")
}

func TestPostNotifications(t *testing.T) {
	router := SetupTestRouter()

	notification := Notification{
		Message: "You have been outbid",
		UserId:  "user1",
	}

	notificationJSON, _ := json.Marshal(notification)

	req, _ := http.NewRequest("POST", "/notifications", bytes.NewBuffer(notificationJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
	assert.Contains(t, w.Body.String(), "Not yet implemented")
}

func TestPostUsersLogin(t *testing.T) {
	router := SetupTestRouter()

	login := Login{
		Username: "testuser",
		Password: "password123",
	}

	loginJSON, _ := json.Marshal(login)

	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
	assert.Contains(t, w.Body.String(), "Not yet implemented")
}

func TestPostUsersRegister(t *testing.T) {
	router := SetupTestRouter()

	user := User{
		Username: "testuser",
		Password: "password123",
	}

	userJSON, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
	assert.Contains(t, w.Body.String(), "Not yet implemented")
}

// HELPER FUNCTIONS

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
