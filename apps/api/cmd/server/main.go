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

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	addressH := handler.NewAddressHandler(addressSvc)
	shareH := handler.NewShareHandler(shareSvc, cfg.BaseURL)

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

	// Public routes
	r.Post("/api/v1/auth/register", authH.Register)
	r.Post("/api/v1/auth/login", authH.Login)

	// Token resolution — public
	r.Get("/api/v1/resolve/{token}", shareH.Resolve)
	r.Get("/api/v1/qr/{token}", shareH.QRCode)

	// Protected routes
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

	fmt.Printf("AddrPass API running on :%s\n", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

func findMigrationsDir() string {
	// Try common locations
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
