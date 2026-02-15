package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/saigenix/bidding-system/internal/domain"
)

type BidRepository struct {
	pool *pgxpool.Pool
}

func NewBidRepository(pool *pgxpool.Pool) *BidRepository {
	return &BidRepository{pool: pool}
}

func (r *BidRepository) Create(ctx context.Context, bid *domain.Bid) error {
	query := `
		INSERT INTO bids (id, auction_id, user_id, amount, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, bid.ID, bid.AuctionID, bid.UserID, bid.Amount, bid.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create bid: %w", err)
	}
	return nil
}

func (r *BidRepository) GetByAuctionID(ctx context.Context, auctionID string) ([]*domain.Bid, error) {
	query := `
		SELECT id, auction_id, user_id, amount, created_at
		FROM bids
		WHERE auction_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, auctionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bids: %w", err)
	}
	defer rows.Close()

	var bids []*domain.Bid
	for rows.Next() {
		var bid domain.Bid
		if err := rows.Scan(&bid.ID, &bid.AuctionID, &bid.UserID, &bid.Amount, &bid.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan bid: %w", err)
		}
		bids = append(bids, &bid)
	}

	return bids, nil
}

func (r *BidRepository) GetHighestBid(ctx context.Context, auctionID string) (*domain.Bid, error) {
	query := `
		SELECT id, auction_id, user_id, amount, created_at
		FROM bids
		WHERE auction_id = $1
		ORDER BY amount DESC, created_at ASC
		LIMIT 1
	`
	var bid domain.Bid
	err := r.pool.QueryRow(ctx, query, auctionID).Scan(
		&bid.ID, &bid.AuctionID, &bid.UserID, &bid.Amount, &bid.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil // No bids yet
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get highest bid: %w", err)
	}
	return &bid, nil
}
