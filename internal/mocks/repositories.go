package mocks

import (
	"context"
	"fmt"
	"sync"

	"github.com/saigenix/bidding-system/internal/domain"
)

// ============================================================================
// MockUserRepository
// ============================================================================

type MockUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User // keyed by ID
	err   error                   // if set, all operations return this error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{users: make(map[string]*domain.User)}
}

func (m *MockUserRepository) SetError(err error) {
	m.err = err
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.err != nil {
		return m.err
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check for duplicate email
	for _, u := range m.users {
		if u.Email == user.Email {
			return fmt.Errorf("user with email %s already exists", user.Email)
		}
	}

	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, fmt.Errorf("user not found")
}

// ============================================================================
// MockProductRepository
// ============================================================================

type MockProductRepository struct {
	mu       sync.RWMutex
	products map[string]*domain.Product
	err      error
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{products: make(map[string]*domain.Product)}
}

func (m *MockProductRepository) SetError(err error) {
	m.err = err
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	if m.err != nil {
		return m.err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.products[product.ID] = product
	return nil
}

func (m *MockProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	if p, ok := m.products[id]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("product not found")
}

func (m *MockProductRepository) List(ctx context.Context) ([]*domain.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*domain.Product
	for _, p := range m.products {
		result = append(result, p)
	}
	return result, nil
}

// ============================================================================
// MockAuctionRepository
// ============================================================================

type MockAuctionRepository struct {
	mu       sync.RWMutex
	auctions map[string]*domain.Auction
	err      error
}

func NewMockAuctionRepository() *MockAuctionRepository {
	return &MockAuctionRepository{auctions: make(map[string]*domain.Auction)}
}

func (m *MockAuctionRepository) SetError(err error) {
	m.err = err
}

func (m *MockAuctionRepository) Create(ctx context.Context, auction *domain.Auction) error {
	if m.err != nil {
		return m.err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.auctions[auction.ID] = auction
	return nil
}

func (m *MockAuctionRepository) GetByID(ctx context.Context, id string) (*domain.Auction, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	if a, ok := m.auctions[id]; ok {
		return a, nil
	}
	return nil, fmt.Errorf("auction not found")
}

func (m *MockAuctionRepository) List(ctx context.Context) ([]*domain.Auction, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*domain.Auction
	for _, a := range m.auctions {
		result = append(result, a)
	}
	return result, nil
}

func (m *MockAuctionRepository) Update(ctx context.Context, auction *domain.Auction) error {
	if m.err != nil {
		return m.err
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.auctions[auction.ID]; !ok {
		return fmt.Errorf("auction not found")
	}
	m.auctions[auction.ID] = auction
	return nil
}

// ============================================================================
// MockBidRepository
// ============================================================================

type MockBidRepository struct {
	mu   sync.RWMutex
	bids []*domain.Bid
	err  error
}

func NewMockBidRepository() *MockBidRepository {
	return &MockBidRepository{}
}

func (m *MockBidRepository) SetError(err error) {
	m.err = err
}

func (m *MockBidRepository) Create(ctx context.Context, bid *domain.Bid) error {
	if m.err != nil {
		return m.err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.bids = append(m.bids, bid)
	return nil
}

func (m *MockBidRepository) GetByAuctionID(ctx context.Context, auctionID string) ([]*domain.Bid, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*domain.Bid
	for _, b := range m.bids {
		if b.AuctionID == auctionID {
			result = append(result, b)
		}
	}
	return result, nil
}

func (m *MockBidRepository) GetHighestBid(ctx context.Context, auctionID string) (*domain.Bid, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	var highest *domain.Bid
	for _, b := range m.bids {
		if b.AuctionID == auctionID {
			if highest == nil || b.Amount > highest.Amount {
				highest = b
			}
		}
	}
	if highest == nil {
		return nil, fmt.Errorf("no bids found")
	}
	return highest, nil
}
