package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"myapp/internal/config"
	"myapp/internal/handler"
	"myapp/internal/repository"
	"myapp/internal/service"
)

type App struct {
	cfg         *config.Config
	httpServer  *http.Server
	mongoClient *mongo.Client
}

func New(cfg *config.Config) (*App, error) {
	app := &App{cfg: cfg}

	// Initialize MongoDB
	mongoClient, err := initMongoDB(cfg)
	if err != nil {
		return nil, err
	}
	app.mongoClient = mongoClient

	// Initialize dependencies
	repos, err := repository.NewRepositories(mongoClient.Database(cfg.MongoDB.Database))
	if err != nil {
		return nil, err
	}
	services := service.NewServices(repos)
	handlers := handler.NewHandlers(services)

	// Setup router
	router := SetupRouter(handlers)
	app.initHTTPServer(router)

	return app, nil
}

func initMongoDB(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoDB.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Verify the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return client, nil
}

func (a *App) initHTTPServer(router *chi.Mux) {
	a.httpServer = &http.Server{
		Addr:         ":" + a.cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  a.cfg.Server.ReadTimeout,
		WriteTimeout: a.cfg.Server.WriteTimeout,
		IdleTimeout:  time.Minute,
	}
}

func (a *App) RunOld(ctx context.Context) error {
	go func() {
		fmt.Printf("Starting server on port %s\n", a.cfg.Server.Port)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return a.Shutdown(shutdownCtx)
}

func (a *App) Run(ctx context.Context) error {
	// Create error channel for server errors
	serverErr := make(chan error, 1)

	go func() {
		fmt.Printf("Starting server on port %s\n", a.cfg.Server.Port)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for either context cancellation or server error
	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	case err := <-serverErr:
		log.Printf("Server error: %v\n", err)
		return err
	}

	// Create shutdown context
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v\n", err)
		return err
	}

	log.Println("Server shut down gracefully")
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	if a.mongoClient != nil {
		return a.mongoClient.Disconnect(ctx)
	}
	return nil
}
