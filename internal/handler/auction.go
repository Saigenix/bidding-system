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
	ProductID     string    `json:"product_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	StartTime     time.Time `json:"start_time" binding:"required" example:"2026-03-01T10:00:00Z"`
	EndTime       time.Time `json:"end_time" binding:"required" example:"2026-03-02T10:00:00Z"`
	StartingPrice float64   `json:"starting_price" binding:"required,min=0" example:"100.00"`
}

// Create godoc
// @Summary      Create an auction
// @Description  Create a new auction for a product with a time window and starting price
// @Tags         Auctions
// @Accept       json
// @Produce      json
// @Param        request  body      CreateAuctionRequest  true  "Auction details"
// @Success      201      {object}  domain.Auction
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /auctions [post]
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

// Get godoc
// @Summary      Get an auction
// @Description  Get an auction by its ID
// @Tags         Auctions
// @Produce      json
// @Param        id   path      string  true  "Auction ID"
// @Success      200  {object}  domain.Auction
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /auctions/{id} [get]
func (h *AuctionHandler) Get(c *gin.Context) {
	id := c.Param("id")
	auction, err := h.auctionService.GetAuction(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "auction not found"})
		return
	}

	c.JSON(http.StatusOK, auction)
}

// List godoc
// @Summary      List all auctions
// @Description  Get a list of all auctions
// @Tags         Auctions
// @Produce      json
// @Success      200  {array}   domain.Auction
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /auctions [get]
func (h *AuctionHandler) List(c *gin.Context) {
	auctions, err := h.auctionService.ListAuctions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list auctions"})
		return
	}

	c.JSON(http.StatusOK, auctions)
}

// Start godoc
// @Summary      Start an auction
// @Description  Transition an auction from pending to active status
// @Tags         Auctions
// @Produce      json
// @Param        id   path      string  true  "Auction ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /auctions/{id}/start [post]
func (h *AuctionHandler) Start(c *gin.Context) {
	id := c.Param("id")
	if err := h.auctionService.StartAuction(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "auction started"})
}

// End godoc
// @Summary      End an auction
// @Description  Transition an auction to ended status, no more bids accepted
// @Tags         Auctions
// @Produce      json
// @Param        id   path      string  true  "Auction ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /auctions/{id}/end [post]
func (h *AuctionHandler) End(c *gin.Context) {
	id := c.Param("id")
	if err := h.auctionService.EndAuction(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "auction ended"})
}
