package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/saigenix/bidding-system/internal/domain"
)

type AuctionService struct {
	auctionRepo domain.AuctionRepository
}

func NewAuctionService(auctionRepo domain.AuctionRepository) *AuctionService {
	return &AuctionService{auctionRepo: auctionRepo}
}

func (s *AuctionService) CreateAuction(ctx context.Context, productID string, startTime, endTime time.Time, startingPrice float64) (*domain.Auction, error) {
	if endTime.Before(startTime) {
		return nil, fmt.Errorf("end time must be after start time")
	}
	if startingPrice < 0 {
		return nil, fmt.Errorf("starting price must be non-negative")
	}

	auction := &domain.Auction{
		ID:            uuid.New().String(),
		ProductID:     productID,
		StartTime:     startTime,
		EndTime:       endTime,
		StartingPrice: startingPrice,
		CurrentPrice:  startingPrice,
		Status:        domain.AuctionStatusPending,
		CreatedAt:     time.Now(),
	}

	if err := s.auctionRepo.Create(ctx, auction); err != nil {
		return nil, fmt.Errorf("failed to create auction: %w", err)
	}

	return auction, nil
}

func (s *AuctionService) GetAuction(ctx context.Context, id string) (*domain.Auction, error) {
	auction, err := s.auctionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get auction: %w", err)
	}
	return auction, nil
}

func (s *AuctionService) ListAuctions(ctx context.Context) ([]*domain.Auction, error) {
	auctions, err := s.auctionRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list auctions: %w", err)
	}
	return auctions, nil
}

func (s *AuctionService) StartAuction(ctx context.Context, id string) error {
	auction, err := s.auctionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get auction: %w", err)
	}

	if auction.Status != domain.AuctionStatusPending {
		return fmt.Errorf("can only start pending auctions")
	}

	auction.Status = domain.AuctionStatusActive
	if err := s.auctionRepo.Update(ctx, auction); err != nil {
		return fmt.Errorf("failed to update auction: %w", err)
	}

	return nil
}

func (s *AuctionService) EndAuction(ctx context.Context, id string) error {
	auction, err := s.auctionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get auction: %w", err)
	}

	if auction.Status == domain.AuctionStatusEnded {
		return fmt.Errorf("auction already ended")
	}

	auction.Status = domain.AuctionStatusEnded
	if err := s.auctionRepo.Update(ctx, auction); err != nil {
		return fmt.Errorf("failed to update auction: %w", err)
	}

	return nil
}
