package domain

import (
	"time"
)

// User represents a user in the bidding system
type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
