package domain

import (
	"time"
)

// Product represents a product that can be auctioned
type Product struct {
	ID          string
	Name        string
	Description string
	OwnerID     string
	CreatedAt   time.Time
}
