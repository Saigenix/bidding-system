package domain

import (
	"time"
)

// Bid represents a bid placed on an auction
type Bid struct {
	ID        string
	AuctionID string
	UserID    string
	Amount    float64
	CreatedAt time.Time
}
