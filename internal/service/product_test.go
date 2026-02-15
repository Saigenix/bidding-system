package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/saigenix/bidding-system/internal/mocks"
)

func newTestProductService() (*ProductService, *mocks.MockProductRepository) {
	repo := mocks.NewMockProductRepository()
	svc := NewProductService(repo)
	return svc, repo
}

// ============================================================================
// CreateProduct
// ============================================================================

func TestProductService_CreateProduct_Success(t *testing.T) {
	svc, _ := newTestProductService()

	product, err := svc.CreateProduct(context.Background(), "Laptop", "Gaming laptop", "owner-123")
	if err != nil {
		t.Fatalf("CreateProduct() unexpected error: %v", err)
	}
	if product == nil {
		t.Fatal("CreateProduct() returned nil product")
	}
	if product.Name != "Laptop" {
		t.Errorf("CreateProduct() name = %q, want %q", product.Name, "Laptop")
	}
	if product.Description != "Gaming laptop" {
		t.Errorf("CreateProduct() description = %q, want %q", product.Description, "Gaming laptop")
	}
	if product.OwnerID != "owner-123" {
		t.Errorf("CreateProduct() ownerID = %q, want %q", product.OwnerID, "owner-123")
	}
	if product.ID == "" {
		t.Error("CreateProduct() product ID is empty")
	}
}

func TestProductService_CreateProduct_RepoError(t *testing.T) {
	svc, repo := newTestProductService()
	repo.SetError(fmt.Errorf("database error"))

	_, err := svc.CreateProduct(context.Background(), "Laptop", "Gaming laptop", "owner-123")
	if err == nil {
		t.Error("CreateProduct() expected error, got nil")
	}
}

// ============================================================================
// GetProduct
// ============================================================================

func TestProductService_GetProduct_Success(t *testing.T) {
	svc, _ := newTestProductService()

	created, _ := svc.CreateProduct(context.Background(), "Laptop", "Gaming laptop", "owner-123")

	product, err := svc.GetProduct(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetProduct() unexpected error: %v", err)
	}
	if product.Name != "Laptop" {
		t.Errorf("GetProduct() name = %q, want %q", product.Name, "Laptop")
	}
}

func TestProductService_GetProduct_NotFound(t *testing.T) {
	svc, _ := newTestProductService()

	_, err := svc.GetProduct(context.Background(), "nonexistent-id")
	if err == nil {
		t.Error("GetProduct() expected error for non-existent product, got nil")
	}
}

// ============================================================================
// ListProducts
// ============================================================================

func TestProductService_ListProducts_Empty(t *testing.T) {
	svc, _ := newTestProductService()

	products, err := svc.ListProducts(context.Background())
	if err != nil {
		t.Fatalf("ListProducts() unexpected error: %v", err)
	}
	if len(products) != 0 {
		t.Errorf("ListProducts() returned %d products, want 0", len(products))
	}
}

func TestProductService_ListProducts_Multiple(t *testing.T) {
	svc, _ := newTestProductService()

	svc.CreateProduct(context.Background(), "Laptop", "Gaming laptop", "owner-1")
	svc.CreateProduct(context.Background(), "Phone", "Smartphone", "owner-2")

	products, err := svc.ListProducts(context.Background())
	if err != nil {
		t.Fatalf("ListProducts() unexpected error: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("ListProducts() returned %d products, want 2", len(products))
	}
}
