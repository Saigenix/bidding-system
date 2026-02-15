package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saigenix/bidding-system/pkg/web"
	"github.com/saigenix/bidding-system/sdk"

	_ "github.com/saigenix/bidding-system/docs/swagger"
)

// @title           Bidding System API
// @version         1.0
// @description     A pluggable, real-time bidding system SDK. Provides auction lifecycle management, real-time bid streaming via SSE and WebSocket, and JWT-based authentication.
// @termsOfService  https://github.com/saigenix/bidding-system

// @contact.name   Saigenix
// @contact.url    https://github.com/saigenix/bidding-system
// @contact.email  support@saigenix.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your JWT token as: Bearer <token>

func main() {
	// Create SDK engine
	engine, err := sdk.NewEngine()
	if err != nil {
		fmt.Printf("Failed to create engine: %v\n", err)
		os.Exit(1)
	}

	// Start engine
	if err := engine.Start(); err != nil {
		fmt.Printf("Failed to start engine: %v\n", err)
		os.Exit(1)
	}

	// Setup router with all handlers
	router := web.SetupRouter(
		engine.AuthService,
		engine.ProductService,
		engine.AuctionService,
		engine.BidService,
	)

	// Create HTTP server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		engine.GetLogger().Info().Msgf("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			engine.GetLogger().Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	engine.GetLogger().Info().Msgf("Bidding system server is running on http://localhost:%s", port)
	engine.GetLogger().Info().Msg("Health check: http://localhost:" + port + "/health")
	engine.GetLogger().Info().Msg("Swagger docs: http://localhost:" + port + "/swagger/index.html")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	engine.GetLogger().Info().Msg("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		engine.GetLogger().Fatal().Err(err).Msg("Server forced to shutdown")
	}

	// Stop engine
	if err := engine.Stop(); err != nil {
		engine.GetLogger().Error().Err(err).Msg("Failed to stop engine")
	}

	engine.GetLogger().Info().Msg("Server exited")
}
