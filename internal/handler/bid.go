package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/saigenix/bidding-system/internal/service"
)

type BidHandler struct {
	bidService *service.BidService
	upgrader   websocket.Upgrader
}

func NewBidHandler(bidService *service.BidService) *BidHandler {
	return &BidHandler{
		bidService: bidService,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for demo (restrict in production)
			},
		},
	}
}

type PlaceBidRequest struct {
	AuctionID string  `json:"auction_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,min=0"`
}

func (h *BidHandler) PlaceBid(c *gin.Context) {
	var req PlaceBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	bid, err := h.bidService.PlaceBid(c.Request.Context(), req.AuctionID, userID.(string), req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bid)
}

func (h *BidHandler) GetBids(c *gin.Context) {
	auctionID := c.Param("auction_id")
	bids, err := h.bidService.GetBids(c.Request.Context(), auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get bids"})
		return
	}

	c.JSON(http.StatusOK, bids)
}

// StreamBids handles SSE streaming of bid updates
func (h *BidHandler) StreamBids(c *gin.Context) {
	auctionID := c.Param("auction_id")

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// In a real implementation, you'd use a pub/sub system (Redis, NATS, etc.)
	// For now, we'll poll every 2 seconds (simplified for demo)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			bids, err := h.bidService.GetBids(c.Request.Context(), auctionID)
			if err != nil {
				continue
			}

			if len(bids) > 0 {
				latestBid := bids[0]
				c.SSEvent("bid", gin.H{
					"id":         latestBid.ID,
					"user_id":    latestBid.UserID,
					"amount":     latestBid.Amount,
					"created_at": latestBid.CreatedAt,
				})
				c.Writer.Flush()
			}
		}
	}
}

// WebSocketHandler handles WebSocket connections for real-time bidding
func (h *BidHandler) WebSocketHandler(c *gin.Context) {
	auctionID := c.Param("auction_id")

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Send initial bids
	bids, _ := h.bidService.GetBids(c.Request.Context(), auctionID)
	conn.WriteJSON(gin.H{"type": "initial", "bids": bids})

	// In a real implementation, use pub/sub for real-time updates
	// For demo, we'll poll every 2 seconds
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	done := make(chan struct{})

	// Read messages from client (e.g., new bid placement via WS)
	go func() {
		defer close(done)
		for {
			var msg map[string]interface{}
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}
			// Handle incoming messages (e.g., bid placement)
			fmt.Printf("Received message: %v\n", msg)
		}
	}()

	// Send updates to client
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			bids, err := h.bidService.GetBids(c.Request.Context(), auctionID)
			if err != nil {
				continue
			}
			if len(bids) > 0 {
				conn.WriteJSON(gin.H{"type": "update", "latest_bid": bids[0]})
			}
		}
	}
}
