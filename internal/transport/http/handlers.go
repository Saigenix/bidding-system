package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewServer() Server {
	return Server{}
}

type Server struct{}

func (s Server) GetAuctions(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, "Not yet implemented")
}

func (s Server) PostAuctions(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, "Not yet implemented")
}

func (s Server) DeleteAuctionsAuctionId(c *gin.Context, auctionId string) {
    c.JSON(http.StatusNotImplemented, "Not yet implemented")
}

func (s Server) GetAuctionsAuctionId(c *gin.Context, auctionId string) {
    c.JSON(http.StatusNotImplemented, "Not yet implemented")
}

func (s Server) PutAuctionsAuctionId(c *gin.Context, auctionId string) {
    c.JSON(http.StatusNotImplemented, "Not yet implemented")
}

func (s Server) PostBids(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, "Not yet implemented")
}

func (s Server) GetBidsBidId(c *gin.Context, bidId string) {
    c.JSON(http.StatusNotImplemented, "Not yet implemented")
}

func (s Server) PostNotifications(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, "Not yet implemented")
}

func (s Server) PostUsersLogin(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, "Not yet implemented")
}

func (s Server) PostUsersRegister(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, "Not yet implemented")
}

func (s Server) PostUsers(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, "Not yet implemented")
}


