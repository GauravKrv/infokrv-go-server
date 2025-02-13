package app

import (
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    
    customMiddleware "myapp/internal/middleware"
    "myapp/internal/handler"
)

func SetupRouter(handlers *handler.Handlers) *chi.Mux {
    router := chi.NewRouter()

    // Middleware
    router.Use(middleware.RequestID)
    router.Use(middleware.RealIP)
    router.Use(middleware.Logger)
    router.Use(middleware.Recoverer)
    router.Use(middleware.Timeout(60 * time.Second))
    router.Use(customMiddleware.Cors)

    // API Routes
    router.Route("/api/v1", func(r chi.Router) {
        r.Route("/users", func(r chi.Router) {
            r.Post("/", handlers.User.Create)
            r.Get("/", handlers.User.List)
            r.Route("/{id}", func(r chi.Router) {
                r.Get("/", handlers.User.Get)
                r.Put("/", handlers.User.Update)
                r.Delete("/", handlers.User.Delete)
            })
        })

        r.Route("/sectionDetail", func(r chi.Router) {
            r.Post("/", handlers.SectionDetail.Create)
            r.Get("/", handlers.SectionDetail.List)
            r.Route("/{id}", func(r chi.Router) {
                r.Get("/", handlers.SectionDetail.Get)
                r.Put("/", handlers.SectionDetail.Update)
                r.Delete("/", handlers.SectionDetail.Delete)
            })
            r.Get("/sectionType/{sectionType}", handlers.SectionDetail.FindBySectionType)
        })
    })

    // Health Check
    router.Get("/health", healthCheckHandler)

    return router
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

