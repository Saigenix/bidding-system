package service

import (
	"context"
	"testing"
	"time"

	"github.com/saigenix/bidding-system/internal/domain"
	"github.com/saigenix/bidding-system/internal/mocks"
)

func newTestBidService() (*BidService, *mocks.MockBidRepository, *mocks.MockAuctionRepository) {
	bidRepo := mocks.NewMockBidRepository()
	auctionRepo := mocks.NewMockAuctionRepository()
	svc := NewBidService(bidRepo, auctionRepo)
	return svc, bidRepo, auctionRepo
}

// createActiveAuction creates an active auction in the mock repo for testing
func createActiveAuction(t *testing.T, auctionRepo *mocks.MockAuctionRepository) *domain.Auction {
	t.Helper()
	auction := &domain.Auction{
		ID:            "auction-123",
		ProductID:     "product-123",
		StartTime:     time.Now().Add(-1 * time.Hour),
		EndTime:       time.Now().Add(1 * time.Hour),
		StartingPrice: 100.00,
		CurrentPrice:  100.00,
		Status:        domain.AuctionStatusActive,
		CreatedAt:     time.Now(),
	}
	if err := auctionRepo.Create(context.Background(), auction); err != nil {
		t.Fatalf("Failed to create test auction: %v", err)
	}
	return auction
}

// ============================================================================
// PlaceBid
// ============================================================================

func TestBidService_PlaceBid_Success(t *testing.T) {
	svc, _, auctionRepo := newTestBidService()
	createActiveAuction(t, auctionRepo)

	bid, err := svc.PlaceBid(context.Background(), "auction-123", "user-456", 150.00)
	if err != nil {
		t.Fatalf("PlaceBid() unexpected error: %v", err)
	}
	if bid == nil {
		t.Fatal("PlaceBid() returned nil bid")
	}
	if bid.AuctionID != "auction-123" {
		t.Errorf("PlaceBid() auctionID = %q, want %q", bid.AuctionID, "auction-123")
	}
	if bid.UserID != "user-456" {
		t.Errorf("PlaceBid() userID = %q, want %q", bid.UserID, "user-456")
	}
	if bid.Amount != 150.00 {
		t.Errorf("PlaceBid() amount = %f, want %f", bid.Amount, 150.00)
	}

	// Verify auction current price was updated
	auction, _ := auctionRepo.GetByID(context.Background(), "auction-123")
	if auction.CurrentPrice != 150.00 {
		t.Errorf("Auction current price = %f, want %f", auction.CurrentPrice, 150.00)
	}
}

func TestBidService_PlaceBid_InactiveAuction(t *testing.T) {
	svc, _, auctionRepo := newTestBidService()

	// Create a pending (not active) auction
	auction := &domain.Auction{
		ID:            "auction-pending",
		ProductID:     "product-123",
		StartTime:     time.Now().Add(1 * time.Hour),
		EndTime:       time.Now().Add(24 * time.Hour),
		StartingPrice: 100.00,
		CurrentPrice:  100.00,
		Status:        domain.AuctionStatusPending,
		CreatedAt:     time.Now(),
	}
	auctionRepo.Create(context.Background(), auction)

	_, err := svc.PlaceBid(context.Background(), "auction-pending", "user-456", 150.00)
	if err == nil {
		t.Error("PlaceBid() expected error for inactive auction, got nil")
	}
}

func TestBidService_PlaceBid_AmountTooLow(t *testing.T) {
	svc, _, auctionRepo := newTestBidService()
	createActiveAuction(t, auctionRepo) // current price is 100.00

	_, err := svc.PlaceBid(context.Background(), "auction-123", "user-456", 50.00)
	if err == nil {
		t.Error("PlaceBid() expected error for bid lower than current price, got nil")
	}
}

func TestBidService_PlaceBid_AmountEqualToCurrentPrice(t *testing.T) {
	svc, _, auctionRepo := newTestBidService()
	createActiveAuction(t, auctionRepo) // current price is 100.00

	_, err := svc.PlaceBid(context.Background(), "auction-123", "user-456", 100.00)
	if err == nil {
		t.Error("PlaceBid() expected error for bid equal to current price, got nil")
	}
}

func TestBidService_PlaceBid_NonExistentAuction(t *testing.T) {
	svc, _, _ := newTestBidService()

	_, err := svc.PlaceBid(context.Background(), "nonexistent", "user-456", 150.00)
	if err == nil {
		t.Error("PlaceBid() expected error for non-existent auction, got nil")
	}
}

// ============================================================================
// GetBids
// ============================================================================

func TestBidService_GetBids_Success(t *testing.T) {
	svc, _, auctionRepo := newTestBidService()
	createActiveAuction(t, auctionRepo)

	svc.PlaceBid(context.Background(), "auction-123", "user-1", 150.00)
	svc.PlaceBid(context.Background(), "auction-123", "user-2", 200.00)

	bids, err := svc.GetBids(context.Background(), "auction-123")
	if err != nil {
		t.Fatalf("GetBids() unexpected error: %v", err)
	}
	if len(bids) != 2 {
		t.Errorf("GetBids() returned %d bids, want 2", len(bids))
	}
}

func TestBidService_GetBids_EmptyAuction(t *testing.T) {
	svc, _, _ := newTestBidService()

	bids, err := svc.GetBids(context.Background(), "auction-no-bids")
	if err != nil {
		t.Fatalf("GetBids() unexpected error: %v", err)
	}
	if len(bids) != 0 {
		t.Errorf("GetBids() returned %d bids, want 0", len(bids))
	}
}

// ============================================================================
// GetWinningBid
// ============================================================================

func TestBidService_GetWinningBid_Success(t *testing.T) {
	svc, _, auctionRepo := newTestBidService()
	createActiveAuction(t, auctionRepo)

	svc.PlaceBid(context.Background(), "auction-123", "user-1", 150.00)
	svc.PlaceBid(context.Background(), "auction-123", "user-2", 200.00)
	svc.PlaceBid(context.Background(), "auction-123", "user-3", 250.00)

	winner, err := svc.GetWinningBid(context.Background(), "auction-123")
	if err != nil {
		t.Fatalf("GetWinningBid() unexpected error: %v", err)
	}
	if winner.Amount != 250.00 {
		t.Errorf("GetWinningBid() amount = %f, want %f", winner.Amount, 250.00)
	}
	if winner.UserID != "user-3" {
		t.Errorf("GetWinningBid() userID = %q, want %q", winner.UserID, "user-3")
	}
}

func TestBidService_GetWinningBid_NoBids(t *testing.T) {
	svc, _, _ := newTestBidService()

	_, err := svc.GetWinningBid(context.Background(), "auction-no-bids")
	if err == nil {
		t.Error("GetWinningBid() expected error when no bids exist, got nil")
	}
}
