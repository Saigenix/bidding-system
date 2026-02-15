package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/saigenix/bidding-system/internal/domain"
)

type BidService struct {
	bidRepo     domain.BidRepository
	auctionRepo domain.AuctionRepository
}

func NewBidService(bidRepo domain.BidRepository, auctionRepo domain.AuctionRepository) *BidService {
	return &BidService{
		bidRepo:     bidRepo,
		auctionRepo: auctionRepo,
	}
}

func (s *BidService) PlaceBid(ctx context.Context, auctionID, userID string, amount float64) (*domain.Bid, error) {
	// Get auction
	auction, err := s.auctionRepo.GetByID(ctx, auctionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get auction: %w", err)
	}

	// Validate auction is active
	if !auction.IsActive() {
		return nil, fmt.Errorf("auction is not active")
	}

	// Validate bid amount is higher than current price
	if amount <= auction.CurrentPrice {
		return nil, fmt.Errorf("bid amount must be higher than current price (%.2f)", auction.CurrentPrice)
	}

	// Create bid
	bid := &domain.Bid{
		ID:        uuid.New().String(),
		AuctionID: auctionID,
		UserID:    userID,
		Amount:    amount,
		CreatedAt: time.Now(),
	}

	if err := s.bidRepo.Create(ctx, bid); err != nil {
		return nil, fmt.Errorf("failed to create bid: %w", err)
	}

	// Update auction current price
	auction.CurrentPrice = amount
	if err := s.auctionRepo.Update(ctx, auction); err != nil {
		return nil, fmt.Errorf("failed to update auction: %w", err)
	}

	return bid, nil
}

func (s *BidService) GetBids(ctx context.Context, auctionID string) ([]*domain.Bid, error) {
	bids, err := s.bidRepo.GetByAuctionID(ctx, auctionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bids: %w", err)
	}
	return bids, nil
}

func (s *BidService) GetWinningBid(ctx context.Context, auctionID string) (*domain.Bid, error) {
	bid, err := s.bidRepo.GetHighestBid(ctx, auctionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get winning bid: %w", err)
	}
	return bid, nil
}
