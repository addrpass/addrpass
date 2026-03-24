package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/addrpass/addrpass/apps/api/internal/config"
	"github.com/addrpass/addrpass/apps/api/internal/database"
	"github.com/addrpass/addrpass/apps/api/internal/handler"
	"github.com/addrpass/addrpass/apps/api/internal/middleware"
	"github.com/addrpass/addrpass/apps/api/internal/service"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Database
	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Run migrations
	if err := database.RunMigrations(ctx, pool, findMigrationsDir()); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Services
	authSvc := service.NewAuthService(pool, cfg.JWTSecret)
	addressSvc := service.NewAddressService(pool)
	shareSvc := service.NewShareService(pool)
	businessSvc := service.NewBusinessService(pool, cfg.JWTSecret)
	webhookSvc := service.NewWebhookService(pool)

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	addressH := handler.NewAddressHandler(addressSvc)
	shareH := handler.NewShareHandler(shareSvc, webhookSvc, cfg.BaseURL, cfg.JWTSecret)
	businessH := handler.NewBusinessHandler(businessSvc)
	webhookH := handler.NewWebhookHandler(webhookSvc)

	// Rate limiters
	publicLimiter := middleware.NewRateLimiter(60, time.Minute)    // 60 req/min for public endpoints
	authLimiter := middleware.NewRateLimiter(10, time.Minute)      // 10 req/min for login/register
	resolveLimiter := middleware.NewRateLimiter(120, time.Minute)  // 120 req/min for token resolution

	// Router
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Public routes (rate limited)
	r.Group(func(r chi.Router) {
		r.Use(authLimiter.Handler)
		r.Post("/api/v1/auth/register", authH.Register)
		r.Post("/api/v1/auth/login", authH.Login)
	})

	// OAuth token endpoint (rate limited)
	r.Group(func(r chi.Router) {
		r.Use(authLimiter.Handler)
		r.Post("/api/v1/oauth/token", businessH.OAuthToken)
	})

	// Token resolution — public but rate limited
	r.Group(func(r chi.Router) {
		r.Use(resolveLimiter.Handler)
		r.Get("/api/v1/resolve/{token}", shareH.Resolve)
	})

	// QR code — public, cached
	r.Group(func(r chi.Router) {
		r.Use(publicLimiter.Handler)
		r.Get("/api/v1/qr/{token}", shareH.QRCode)
	})

	// Protected routes (user JWT)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWTSecret))

		// User
		r.Get("/api/v1/auth/me", authH.Me)

		// Addresses
		r.Post("/api/v1/addresses", addressH.Create)
		r.Get("/api/v1/addresses", addressH.List)
		r.Get("/api/v1/addresses/{id}", addressH.Get)
		r.Put("/api/v1/addresses/{id}", addressH.Update)
		r.Delete("/api/v1/addresses/{id}", addressH.Delete)

		// Shares
		r.Post("/api/v1/shares", shareH.Create)
		r.Get("/api/v1/shares", shareH.List)
		r.Patch("/api/v1/shares/{id}/revoke", shareH.Revoke)
		r.Delete("/api/v1/shares/{id}", shareH.Delete)
		r.Get("/api/v1/shares/{id}/accesses", shareH.AccessLogs)

		// Businesses
		r.Post("/api/v1/businesses", businessH.CreateBusiness)
		r.Get("/api/v1/businesses", businessH.ListBusinesses)
		r.Post("/api/v1/businesses/{businessId}/api-keys", businessH.CreateAPIKey)
		r.Get("/api/v1/businesses/{businessId}/api-keys", businessH.ListAPIKeys)
		r.Patch("/api/v1/api-keys/{keyId}/revoke", businessH.RevokeAPIKey)

		// Webhooks
		r.Post("/api/v1/webhooks", webhookH.Create)
		r.Get("/api/v1/webhooks", webhookH.List)
		r.Delete("/api/v1/webhooks/{id}", webhookH.Delete)
	})

	// Server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		log.Println("Shutting down server...")
		srv.Shutdown(shutdownCtx)
	}()

	fmt.Printf("AddrPass API v2 running on :%s\n", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

func findMigrationsDir() string {
	dirs := []string{
		"migrations",
		"apps/api/migrations",
		"/app/migrations",
	}
	for _, d := range dirs {
		if _, err := os.Stat(d); err == nil {
			return d
		}
	}
	return "migrations"
}
