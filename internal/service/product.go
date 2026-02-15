package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/saigenix/bidding-system/internal/domain"
)

type ProductService struct {
	productRepo domain.ProductRepository
}

func NewProductService(productRepo domain.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) CreateProduct(ctx context.Context, name, description, ownerID string) (*domain.Product, error) {
	product := &domain.Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

func (s *ProductService) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	products, err := s.productRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return products, nil
}
