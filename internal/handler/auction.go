package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saigenix/bidding-system/internal/service"
)

type AuctionHandler struct {
	auctionService *service.AuctionService
}

func NewAuctionHandler(auctionService *service.AuctionService) *AuctionHandler {
	return &AuctionHandler{auctionService: auctionService}
}

type CreateAuctionRequest struct {
	ProductID     string    `json:"product_id" binding:"required"`
	StartTime     time.Time `json:"start_time" binding:"required"`
	EndTime       time.Time `json:"end_time" binding:"required"`
	StartingPrice float64   `json:"starting_price" binding:"required,min=0"`
}

func (h *AuctionHandler) Create(c *gin.Context) {
	var req CreateAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	auction, err := h.auctionService.CreateAuction(c.Request.Context(), req.ProductID, req.StartTime, req.EndTime, req.StartingPrice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, auction)
}

func (h *AuctionHandler) Get(c *gin.Context) {
	id := c.Param("id")
	auction, err := h.auctionService.GetAuction(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "auction not found"})
		return
	}

	c.JSON(http.StatusOK, auction)
}

func (h *AuctionHandler) List(c *gin.Context) {
	auctions, err := h.auctionService.ListAuctions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list auctions"})
		return
	}

	c.JSON(http.StatusOK, auctions)
}

func (h *AuctionHandler) Start(c *gin.Context) {
	id := c.Param("id")
	if err := h.auctionService.StartAuction(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "auction started"})
}

func (h *AuctionHandler) End(c *gin.Context) {
	id := c.Param("id")
	if err := h.auctionService.EndAuction(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "auction ended"})
}
