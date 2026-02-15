package sdk

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/saigenix/bidding-system/config"
)

// Option is a functional option for configuring the Engine
type Option func(*Engine) error

// WithConfig sets a custom configuration
func WithConfig(cfg *config.Config) Option {
	return func(e *Engine) error {
		e.cfg = cfg
		return nil
	}
}

// WithLogger sets a custom logger
func WithLogger(logger zerolog.Logger) Option {
	return func(e *Engine) error {
		e.logger = logger
		return nil
	}
}

// WithDBPool sets a custom database pool
func WithDBPool(pool *pgxpool.Pool) Option {
	return func(e *Engine) error {
		e.dbPool = pool
		return nil
	}
}

// WithJWTSecret sets a custom JWT secret
func WithJWTSecret(secret string) Option {
	return func(e *Engine) error {
		if e.cfg == nil {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			e.cfg = cfg
		}
		e.cfg.JWT.Secret = secret
		return nil
	}
}
