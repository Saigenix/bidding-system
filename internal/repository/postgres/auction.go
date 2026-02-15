package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/saigenix/bidding-system/internal/domain"
)

type AuctionRepository struct {
	pool *pgxpool.Pool
}

func NewAuctionRepository(pool *pgxpool.Pool) *AuctionRepository {
	return &AuctionRepository{pool: pool}
}

func (r *AuctionRepository) Create(ctx context.Context, auction *domain.Auction) error {
	query := `
		INSERT INTO auctions (id, product_id, start_time, end_time, starting_price, current_price, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query,
		auction.ID, auction.ProductID, auction.StartTime, auction.EndTime,
		auction.StartingPrice, auction.CurrentPrice, auction.Status, auction.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create auction: %w", err)
	}
	return nil
}

func (r *AuctionRepository) GetByID(ctx context.Context, id string) (*domain.Auction, error) {
	query := `
		SELECT id, product_id, start_time, end_time, starting_price, current_price, status, created_at
		FROM auctions
		WHERE id = $1
	`
	var auction domain.Auction
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&auction.ID, &auction.ProductID, &auction.StartTime, &auction.EndTime,
		&auction.StartingPrice, &auction.CurrentPrice, &auction.Status, &auction.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get auction: %w", err)
	}
	return &auction, nil
}

func (r *AuctionRepository) List(ctx context.Context) ([]*domain.Auction, error) {
	query := `
		SELECT id, product_id, start_time, end_time, starting_price, current_price, status, created_at
		FROM auctions
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list auctions: %w", err)
	}
	defer rows.Close()

	var auctions []*domain.Auction
	for rows.Next() {
		var auction domain.Auction
		if err := rows.Scan(
			&auction.ID, &auction.ProductID, &auction.StartTime, &auction.EndTime,
			&auction.StartingPrice, &auction.CurrentPrice, &auction.Status, &auction.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan auction: %w", err)
		}
		auctions = append(auctions, &auction)
	}

	return auctions, nil
}

func (r *AuctionRepository) Update(ctx context.Context, auction *domain.Auction) error {
	query := `
		UPDATE auctions
		SET product_id = $2, start_time = $3, end_time = $4,
		    starting_price = $5, current_price = $6, status = $7
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query,
		auction.ID, auction.ProductID, auction.StartTime, auction.EndTime,
		auction.StartingPrice, auction.CurrentPrice, auction.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to update auction: %w", err)
	}
	return nil
}
