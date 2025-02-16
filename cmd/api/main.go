package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "github.com/joho/godotenv"
    "myapp/internal/app"
    "myapp/internal/config"
)

func main() {
    // Load .env file if it exists
    if err := godotenv.Load(); err != nil {
        log.Printf("No .env file found or error loading it: %v", err)
    }
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Override port with environment variable if provided
    if port := os.Getenv("PORT"); port != "" {
        cfg.Server.Port = port
    }

    // Create app instance
    application, err := app.New(cfg)
    if err != nil {
        log.Fatalf("Failed to create app: %v", err)
    }

    // Use signal.NotifyContext instead of manual signal handling
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    if err := application.Run(ctx); err != nil {
        log.Fatalf("Failed to run app: %v", err)
    }
}

func mainOld() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Create app instance
    application, err := app.New(cfg)
    if err != nil {
        log.Fatalf("Failed to create app: %v", err)
    }

    // Graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan
        cancel()
    }()

    if err := application.Run(ctx); err != nil {
        log.Fatalf("Failed to run app: %v", err)
    }
}