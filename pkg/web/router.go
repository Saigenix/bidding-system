package web

import (
	"github.com/gin-gonic/gin"
	"github.com/saigenix/bidding-system/internal/auth"
	"github.com/saigenix/bidding-system/internal/handler"
	"github.com/saigenix/bidding-system/internal/service"
)

// SetupRouter initializes and configures the Gin router
func SetupRouter(
	authService *service.AuthService,
	productService *service.ProductService,
	auctionService *service.AuctionService,
	bidService *service.BidService,
) *gin.Engine {
	router := gin.Default()

	// CORS middleware (allow all origins for demo)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	productHandler := handler.NewProductHandler(productService)
	auctionHandler := handler.NewAuctionHandler(auctionService)
	bidHandler := handler.NewBidHandler(bidService)

	// Public routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	// Protected routes (require JWT)
	jwtMiddleware := auth.JWTMiddleware(authService)

	productRoutes := router.Group("/products")
	productRoutes.Use(jwtMiddleware)
	{
		productRoutes.POST("", productHandler.Create)
		productRoutes.GET("/:id", productHandler.Get)
		productRoutes.GET("", productHandler.List)
	}

	auctionRoutes := router.Group("/auctions")
	auctionRoutes.Use(jwtMiddleware)
	{
		auctionRoutes.POST("", auctionHandler.Create)
		auctionRoutes.GET("/:id", auctionHandler.Get)
		auctionRoutes.GET("", auctionHandler.List)
		auctionRoutes.POST("/:id/start", auctionHandler.Start)
		auctionRoutes.POST("/:id/end", auctionHandler.End)

		// Bid routes under auctions
		auctionRoutes.POST("/:auction_id/bids", bidHandler.PlaceBid)
		auctionRoutes.GET("/:auction_id/bids", bidHandler.GetBids)

		// Real-time routes (SSE and WebSocket)
		auctionRoutes.GET("/:auction_id/bids/stream", bidHandler.StreamBids)
		auctionRoutes.GET("/:auction_id/bids/ws", bidHandler.WebSocketHandler)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
