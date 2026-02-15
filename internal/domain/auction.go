package domain

import (
	"time"
)

// AuctionStatus represents the current state of an auction
type AuctionStatus string

const (
	AuctionStatusPending AuctionStatus = "pending"
	AuctionStatusActive  AuctionStatus = "active"
	AuctionStatusEnded   AuctionStatus = "ended"
)

// Auction represents an auction for a product
type Auction struct {
	ID            string
	ProductID     string
	StartTime     time.Time
	EndTime       time.Time
	StartingPrice float64
	CurrentPrice  float64
	Status        AuctionStatus
	CreatedAt     time.Time
}

// IsActive checks if the auction is currently active
func (a *Auction) IsActive() bool {
	now := time.Now()
	return a.Status == AuctionStatusActive && now.After(a.StartTime) && now.Before(a.EndTime)
}

// HasEnded checks if the auction has ended
func (a *Auction) HasEnded() bool {
	return a.Status == AuctionStatusEnded || time.Now().After(a.EndTime)
}
