package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	httpApi "github.com/saigenix/bidding-system/internal/transport/http"

	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware" // Request validator
)

func main() {
	// Create main context which be used in application
	ctx := context.Background()

	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	server := httpApi.NewServer()

	// Do not initialize default gin router, instead initialize a new empty router and adjust configuration as needed.
	//TODO:
	// 1. Adjust logging middleware
	// 2. Adjust security middleware and configuration
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Get spec file. Ensure you build/start application from root path of repository. Either way, spec file wont be found.
	//TODO: Rework this solution to:
	// 1. Load Swagger endpoint from spec
	// 2. Create validator based on swagger endpoint
	// More info here: https://github.com/oapi-codegen/oapi-codegen/blob/main/examples/petstore-expanded/gin/petstore.go
	// 3. Consider moving it to separate middleware.go file which will store all related middlewares.
	curDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %s", err)
	}
	specFile := filepath.Join(curDir, "api", "spec.yaml")

	if v, err := middleware.OapiValidatorFromYamlFile(filepath.Join(specFile)); err != nil {
		log.Fatalf("failed to create validator: %s", err)
	} else {
		r.Use(v)
	}

	// Register the handlers with the main router.
	httpApi.RegisterHandlers(r, server)

	// Initialize the server.
	s := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
	}

	// Run the server in goroutine. This way, it will wait for shutdown notification from function from main code.
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %s\n", err)
		}
		log.Println("Stopping server ...")
	}()

	// As we run server in separate goroutine, this part of code will block and just wait for signals that abort the server. When those are received, we shutdown the server with Shutdown() method. It have timeout to ensure we dont block unepxectedly with that invocation.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctxShutdown, cancelShutdown := context.WithTimeout(ctx, 10*time.Second)
	defer cancelShutdown()

	if err := s.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("HTTP server shutdown error: %s\n", err)
	}
	log.Println("HTTP server stop completed")
}
