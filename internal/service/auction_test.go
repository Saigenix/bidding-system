package service

import (
	"context"
	"testing"
	"time"

	"github.com/saigenix/bidding-system/internal/domain"
	"github.com/saigenix/bidding-system/internal/mocks"
)

func newTestAuctionService() (*AuctionService, *mocks.MockAuctionRepository) {
	repo := mocks.NewMockAuctionRepository()
	svc := NewAuctionService(repo)
	return svc, repo
}

// ============================================================================
// CreateAuction
// ============================================================================

func TestAuctionService_CreateAuction_Success(t *testing.T) {
	svc, _ := newTestAuctionService()

	start := time.Now().Add(1 * time.Hour)
	end := time.Now().Add(24 * time.Hour)

	auction, err := svc.CreateAuction(context.Background(), "product-123", start, end, 100.00)
	if err != nil {
		t.Fatalf("CreateAuction() unexpected error: %v", err)
	}
	if auction == nil {
		t.Fatal("CreateAuction() returned nil auction")
	}
	if auction.ProductID != "product-123" {
		t.Errorf("CreateAuction() productID = %q, want %q", auction.ProductID, "product-123")
	}
	if auction.StartingPrice != 100.00 {
		t.Errorf("CreateAuction() startingPrice = %f, want %f", auction.StartingPrice, 100.00)
	}
	if auction.CurrentPrice != 100.00 {
		t.Errorf("CreateAuction() currentPrice = %f, want %f", auction.CurrentPrice, 100.00)
	}
	if auction.Status != domain.AuctionStatusPending {
		t.Errorf("CreateAuction() status = %q, want %q", auction.Status, domain.AuctionStatusPending)
	}
}

func TestAuctionService_CreateAuction_EndBeforeStart(t *testing.T) {
	svc, _ := newTestAuctionService()

	start := time.Now().Add(24 * time.Hour)
	end := time.Now().Add(1 * time.Hour) // end before start

	_, err := svc.CreateAuction(context.Background(), "product-123", start, end, 100.00)
	if err == nil {
		t.Error("CreateAuction() expected error for end before start, got nil")
	}
}

func TestAuctionService_CreateAuction_NegativePrice(t *testing.T) {
	svc, _ := newTestAuctionService()

	start := time.Now().Add(1 * time.Hour)
	end := time.Now().Add(24 * time.Hour)

	_, err := svc.CreateAuction(context.Background(), "product-123", start, end, -50.00)
	if err == nil {
		t.Error("CreateAuction() expected error for negative price, got nil")
	}
}

// ============================================================================
// GetAuction / ListAuctions
// ============================================================================

func TestAuctionService_GetAuction_Success(t *testing.T) {
	svc, _ := newTestAuctionService()

	start := time.Now().Add(1 * time.Hour)
	end := time.Now().Add(24 * time.Hour)
	created, _ := svc.CreateAuction(context.Background(), "product-123", start, end, 100.00)

	auction, err := svc.GetAuction(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetAuction() unexpected error: %v", err)
	}
	if auction.ID != created.ID {
		t.Errorf("GetAuction() ID = %q, want %q", auction.ID, created.ID)
	}
}

func TestAuctionService_GetAuction_NotFound(t *testing.T) {
	svc, _ := newTestAuctionService()

	_, err := svc.GetAuction(context.Background(), "nonexistent-id")
	if err == nil {
		t.Error("GetAuction() expected error for non-existent auction, got nil")
	}
}

func TestAuctionService_ListAuctions(t *testing.T) {
	svc, _ := newTestAuctionService()

	start := time.Now().Add(1 * time.Hour)
	end := time.Now().Add(24 * time.Hour)
	svc.CreateAuction(context.Background(), "product-1", start, end, 100.00)
	svc.CreateAuction(context.Background(), "product-2", start, end, 200.00)

	auctions, err := svc.ListAuctions(context.Background())
	if err != nil {
		t.Fatalf("ListAuctions() unexpected error: %v", err)
	}
	if len(auctions) != 2 {
		t.Errorf("ListAuctions() returned %d auctions, want 2", len(auctions))
	}
}

// ============================================================================
// StartAuction
// ============================================================================

func TestAuctionService_StartAuction_Success(t *testing.T) {
	svc, _ := newTestAuctionService()

	start := time.Now().Add(1 * time.Hour)
	end := time.Now().Add(24 * time.Hour)
	auction, _ := svc.CreateAuction(context.Background(), "product-123", start, end, 100.00)

	err := svc.StartAuction(context.Background(), auction.ID)
	if err != nil {
		t.Fatalf("StartAuction() unexpected error: %v", err)
	}

	// Verify status changed
	updated, _ := svc.GetAuction(context.Background(), auction.ID)
	if updated.Status != domain.AuctionStatusActive {
		t.Errorf("StartAuction() status = %q, want %q", updated.Status, domain.AuctionStatusActive)
	}
}

func TestAuctionService_StartAuction_AlreadyActive(t *testing.T) {
	svc, _ := newTestAuctionService()

	start := time.Now().Add(1 * time.Hour)
	end := time.Now().Add(24 * time.Hour)
	auction, _ := svc.CreateAuction(context.Background(), "product-123", start, end, 100.00)

	// Start once
	svc.StartAuction(context.Background(), auction.ID)

	// Try to start again
	err := svc.StartAuction(context.Background(), auction.ID)
	if err == nil {
		t.Error("StartAuction() expected error for non-pending auction, got nil")
	}
}

// ============================================================================
// EndAuction
// ============================================================================

func TestAuctionService_EndAuction_Success(t *testing.T) {
	svc, _ := newTestAuctionService()

	start := time.Now().Add(1 * time.Hour)
	end := time.Now().Add(24 * time.Hour)
	auction, _ := svc.CreateAuction(context.Background(), "product-123", start, end, 100.00)
	svc.StartAuction(context.Background(), auction.ID)

	err := svc.EndAuction(context.Background(), auction.ID)
	if err != nil {
		t.Fatalf("EndAuction() unexpected error: %v", err)
	}

	updated, _ := svc.GetAuction(context.Background(), auction.ID)
	if updated.Status != domain.AuctionStatusEnded {
		t.Errorf("EndAuction() status = %q, want %q", updated.Status, domain.AuctionStatusEnded)
	}
}

func TestAuctionService_EndAuction_AlreadyEnded(t *testing.T) {
	svc, _ := newTestAuctionService()

	start := time.Now().Add(1 * time.Hour)
	end := time.Now().Add(24 * time.Hour)
	auction, _ := svc.CreateAuction(context.Background(), "product-123", start, end, 100.00)
	svc.StartAuction(context.Background(), auction.ID)
	svc.EndAuction(context.Background(), auction.ID)

	// Try to end again
	err := svc.EndAuction(context.Background(), auction.ID)
	if err == nil {
		t.Error("EndAuction() expected error for already ended auction, got nil")
	}
}
