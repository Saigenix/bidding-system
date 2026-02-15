package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/saigenix/bidding-system/internal/domain"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (id, name, description, owner_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, product.ID, product.Name, product.Description, product.OwnerID, product.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	query := `
		SELECT id, name, description, owner_id, created_at
		FROM products
		WHERE id = $1
	`
	var product domain.Product
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.OwnerID, &product.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return &product, nil
}

func (r *ProductRepository) List(ctx context.Context) ([]*domain.Product, error) {
	query := `
		SELECT id, name, description, owner_id, created_at
		FROM products
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.OwnerID, &product.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &product)
	}

	return products, nil
}
