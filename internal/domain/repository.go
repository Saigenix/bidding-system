package domain

import (
	"context"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id string) (*Product, error)
	List(ctx context.Context) ([]*Product, error)
}

// AuctionRepository defines the interface for auction data operations
type AuctionRepository interface {
	Create(ctx context.Context, auction *Auction) error
	GetByID(ctx context.Context, id string) (*Auction, error)
	List(ctx context.Context) ([]*Auction, error)
	Update(ctx context.Context, auction *Auction) error
}

// BidRepository defines the interface for bid data operations
type BidRepository interface {
	Create(ctx context.Context, bid *Bid) error
	GetByAuctionID(ctx context.Context, auctionID string) ([]*Bid, error)
	GetHighestBid(ctx context.Context, auctionID string) (*Bid, error)
}
