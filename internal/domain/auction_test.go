package domain

import (
	"testing"
	"time"
)

func TestAuction_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		auction  Auction
		expected bool
	}{
		{
			name: "active auction within time window",
			auction: Auction{
				Status:    AuctionStatusActive,
				StartTime: time.Now().Add(-1 * time.Hour),
				EndTime:   time.Now().Add(1 * time.Hour),
			},
			expected: true,
		},
		{
			name: "pending auction",
			auction: Auction{
				Status:    AuctionStatusPending,
				StartTime: time.Now().Add(-1 * time.Hour),
				EndTime:   time.Now().Add(1 * time.Hour),
			},
			expected: false,
		},
		{
			name: "ended auction",
			auction: Auction{
				Status:    AuctionStatusEnded,
				StartTime: time.Now().Add(-2 * time.Hour),
				EndTime:   time.Now().Add(-1 * time.Hour),
			},
			expected: false,
		},
		{
			name: "active but before start time",
			auction: Auction{
				Status:    AuctionStatusActive,
				StartTime: time.Now().Add(1 * time.Hour),
				EndTime:   time.Now().Add(2 * time.Hour),
			},
			expected: false,
		},
		{
			name: "active but after end time",
			auction: Auction{
				Status:    AuctionStatusActive,
				StartTime: time.Now().Add(-2 * time.Hour),
				EndTime:   time.Now().Add(-1 * time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.auction.IsActive()
			if result != tt.expected {
				t.Errorf("IsActive() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAuction_HasEnded(t *testing.T) {
	tests := []struct {
		name     string
		auction  Auction
		expected bool
	}{
		{
			name: "status is ended",
			auction: Auction{
				Status:  AuctionStatusEnded,
				EndTime: time.Now().Add(1 * time.Hour),
			},
			expected: true,
		},
		{
			name: "past end time",
			auction: Auction{
				Status:  AuctionStatusActive,
				EndTime: time.Now().Add(-1 * time.Hour),
			},
			expected: true,
		},
		{
			name: "active and before end time",
			auction: Auction{
				Status:  AuctionStatusActive,
				EndTime: time.Now().Add(1 * time.Hour),
			},
			expected: false,
		},
		{
			name: "pending and before end time",
			auction: Auction{
				Status:  AuctionStatusPending,
				EndTime: time.Now().Add(1 * time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.auction.HasEnded()
			if result != tt.expected {
				t.Errorf("HasEnded() = %v, want %v", result, tt.expected)
			}
		})
	}
}
