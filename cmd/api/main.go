package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "myapp/internal/app"
    "myapp/internal/config"
)

func main() {
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


// // cmd/api/main.go
// package main

// import (
//     "context"
//     "log"
//     "os"
//     "os/signal"
//     "syscall"

//     "myapp/internal/app"
//     "myapp/internal/config"
// )

// func main() {
//     // Load configuration
//     cfg, err := config.Load("configs/config.yaml")
//     if err != nil {
//         log.Fatalf("Failed to load config: %v", err)
//     }

//     // Initialize application
//     app, err := app.New(cfg)
//     if err != nil {
//         log.Fatalf("Failed to initialize application: %v", err)
//     }

//     // Setup graceful shutdown
//     ctx, cancel := context.WithCancel(context.Background())
//     defer cancel()

//     // Handle shutdown signals
//     sigChan := make(chan os.Signal, 1)
//     signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

//     go func() {
//         sig := <-sigChan
//         log.Printf("Received shutdown signal: %v", sig)
//         cancel()
//     }()

//     // Run the application
//     if err := app.Run(ctx); err != nil {
//         log.Fatalf("Failed to run application: %v", err)
//     }
// }