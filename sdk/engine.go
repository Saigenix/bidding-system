package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/saigenix/bidding-system/config"
	"github.com/saigenix/bidding-system/internal/domain"
	"github.com/saigenix/bidding-system/internal/repository/postgres"
	"github.com/saigenix/bidding-system/internal/service"
	"github.com/saigenix/bidding-system/pkg/db"
	"github.com/saigenix/bidding-system/pkg/logger"
)

// Engine is the main bidding system SDK
type Engine struct {
	// Configuration
	cfg    *config.Config
	logger zerolog.Logger

	// Infrastructure
	dbPool *pgxpool.Pool

	// Repositories
	userRepo    domain.UserRepository
	productRepo domain.ProductRepository
	auctionRepo domain.AuctionRepository
	bidRepo     domain.BidRepository

	// Services
	AuthService    *service.AuthService
	ProductService *service.ProductService
	AuctionService *service.AuctionService
	BidService     *service.BidService
}

// NewEngine creates a new bidding system engine with the given options
func NewEngine(opts ...Option) (*Engine, error) {
	// Load default config
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	engine := &Engine{
		cfg:    cfg,
		logger: logger.NewLogger(cfg.Logger.Level),
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(engine); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	// Initialize database if not provided
	if engine.dbPool == nil {
		pool, err := db.NewPostgresPool(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create database pool: %w", err)
		}
		engine.dbPool = pool
	}

	// Initialize repositories
	engine.userRepo = postgres.NewUserRepository(engine.dbPool)
	engine.productRepo = postgres.NewProductRepository(engine.dbPool)
	engine.auctionRepo = postgres.NewAuctionRepository(engine.dbPool)
	engine.bidRepo = postgres.NewBidRepository(engine.dbPool)

	// Initialize services
	engine.AuthService = service.NewAuthService(engine.userRepo, cfg.JWT.Secret, cfg.JWT.ExpirationHour)
	engine.ProductService = service.NewProductService(engine.productRepo)
	engine.AuctionService = service.NewAuctionService(engine.auctionRepo)
	engine.BidService = service.NewBidService(engine.bidRepo, engine.auctionRepo)

	engine.logger.Info().Msg("Bidding system engine initialized")
	return engine, nil
}

// Start initializes the engine and starts background workers
func (e *Engine) Start() error {
	e.logger.Info().Msg("Starting bidding system engine")

	// TODO: Start background workers (e.g., auto-end expired auctions)
	// go e.auctionEndWorker()

	return nil
}

// Stop gracefully shuts down the engine
func (e *Engine) Stop() error {
	e.logger.Info().Msg("Stopping bidding system engine")

	if e.dbPool != nil {
		e.dbPool.Close()
	}

	return nil
}

// GetLogger returns the logger
func (e *Engine) GetLogger() *zerolog.Logger {
	return &e.logger
}

// HealthCheck verifies that the engine is healthy
func (e *Engine) HealthCheck(ctx context.Context) error {
	return db.HealthCheck(ctx, e.dbPool)
}

// CreateProduct is a convenience method for creating a product
func (e *Engine) CreateProduct(ctx context.Context, name, description, ownerID string) (*domain.Product, error) {
	return e.ProductService.CreateProduct(ctx, name, description, ownerID)
}

// CreateAuction is a convenience method for creating an auction
func (e *Engine) CreateAuction(ctx context.Context, productID string, startTime, endTime time.Time, startingPrice float64) (*domain.Auction, error) {
	return e.AuctionService.CreateAuction(ctx, productID, startTime, endTime, startingPrice)
}

// PlaceBid is a convenience method for placing a bid
func (e *Engine) PlaceBid(ctx context.Context, auctionID, userID string, amount float64) (*domain.Bid, error) {
	return e.BidService.PlaceBid(ctx, auctionID, userID, amount)
}
