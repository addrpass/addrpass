package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	delegationSvc := service.NewDelegationService(pool)
	labelSvc := service.NewLabelService(pool)
	oauthSvc := service.NewOAuthService(pool, cfg.JWTSecret)
	billingSvc := service.NewBillingService(pool)

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	addressH := handler.NewAddressHandler(addressSvc)
	shareH := handler.NewShareHandler(shareSvc, webhookSvc, billingSvc, cfg.BaseURL, cfg.JWTSecret)
	businessH := handler.NewBusinessHandler(businessSvc)
	webhookH := handler.NewWebhookHandler(webhookSvc)
	delegationH := handler.NewDelegationHandler(delegationSvc)
	labelH := handler.NewLabelHandler(labelSvc)
	oauthH := handler.NewOAuthHandler(oauthSvc, addressSvc)
	billingH := handler.NewBillingHandler(billingSvc)

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

	// Public plan info
	r.Get("/api/v1/plans", billingH.GetPlanLimits)

	// Public routes (rate limited)
	r.Group(func(r chi.Router) {
		r.Use(authLimiter.Handler)
		r.Post("/api/v1/auth/register", authH.Register)
		r.Post("/api/v1/auth/login", authH.Login)
	})

	// OAuth endpoints (rate limited)
	r.Group(func(r chi.Router) {
		r.Use(authLimiter.Handler)
		r.Post("/api/v1/oauth/token", func(w http.ResponseWriter, r *http.Request) {
			// Peek at grant_type to route to correct handler
			// We need to handle both client_credentials and authorization_code
			// Read body, check grant_type, then dispatch
			var peek struct{ GrantType string `json:"grant_type"` }
			body, _ := io.ReadAll(r.Body)
			r.Body.Close()
			json.Unmarshal(body, &peek)
			r.Body = io.NopCloser(bytes.NewReader(body))

			if peek.GrantType == "authorization_code" {
				oauthH.Exchange(w, r)
			} else {
				businessH.OAuthToken(w, r)
			}
		})
	})

	// Token resolution — public but rate limited
	r.Group(func(r chi.Router) {
		r.Use(resolveLimiter.Handler)
		r.Get("/api/v1/resolve/{token}", shareH.Resolve)
	})

	// QR code and labels — public, cached
	r.Group(func(r chi.Router) {
		r.Use(publicLimiter.Handler)
		r.Get("/api/v1/qr/{token}", shareH.QRCode)
		r.Get("/api/v1/labels/{ref}/image", labelH.GetLabelImage)
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

		// Delegations
		r.Post("/api/v1/delegations", delegationH.CreateByUser)
		r.Get("/api/v1/shares/{shareId}/delegations", delegationH.ListForShare)
		r.Patch("/api/v1/delegations/{id}/revoke", delegationH.Revoke)

		// Labels
		r.Post("/api/v1/labels", labelH.Create)

		// Billing & Usage
		r.Get("/api/v1/billing/usage", billingH.GetUsage)

		// OAuth apps
		r.Post("/api/v1/businesses/{businessId}/oauth-apps", oauthH.CreateApp)

		// OAuth consent flow (user must be logged in)
		r.Get("/api/v1/oauth/authorize", oauthH.Authorize)
		r.Post("/api/v1/oauth/consent", oauthH.Consent)
	})

	// Embeddable widget JS (public, cached)
	r.Get("/widget.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte(widgetJS))
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

// Embeddable checkout widget JavaScript
const widgetJS = `(function(){
  "use strict";
  var API="https://api.addrpass.com";
  var APP="https://addrpass.com";

  function AddrPass(opts){
    this.clientId=opts.clientId||"";
    this.redirectUri=opts.redirectUri||window.location.origin+"/addrpass/callback";
    this.scope=opts.scope||"delivery";
    this.onToken=opts.onToken||function(){};
    this.onError=opts.onError||function(){};
  }

  AddrPass.prototype.authorize=function(){
    var state=Math.random().toString(36).substr(2,12);
    var url=APP+"/authorize?client_id="+encodeURIComponent(this.clientId)
      +"&redirect_uri="+encodeURIComponent(this.redirectUri)
      +"&scope="+encodeURIComponent(this.scope)
      +"&state="+state;

    var w=600,h=700;
    var left=(screen.width-w)/2;
    var top=(screen.height-h)/2;
    var popup=window.open(url,"addrpass_consent","width="+w+",height="+h+",left="+left+",top="+top);

    var self=this;
    window.addEventListener("message",function handler(e){
      if(e.origin!==APP)return;
      window.removeEventListener("message",handler);
      if(popup)popup.close();
      if(e.data.error){self.onError(e.data.error);return;}
      self.onToken(e.data);
    });
  };

  AddrPass.prototype.renderButton=function(container){
    var el=typeof container==="string"?document.querySelector(container):container;
    if(!el)return;
    var btn=document.createElement("button");
    btn.type="button";
    btn.innerHTML='<svg width="16" height="16" viewBox="0 0 32 32" fill="none" style="vertical-align:middle;margin-right:6px"><path d="M16 2L4 8v8c0 8.4 5.12 16.24 12 18 6.88-1.76 12-9.6 12-18V8L16 2z" fill="#0F172A"/><circle cx="16" cy="13" r="3" fill="#22D3EE"/><path d="M14 15.5L13 22h6l-1-6.5" fill="#22D3EE" opacity="0.7"/></svg><span>Share via AddrPass</span>';
    btn.style.cssText="display:inline-flex;align-items:center;padding:10px 20px;border-radius:8px;border:1px solid #E2E8F0;background:#fff;color:#0F172A;font-family:-apple-system,sans-serif;font-size:14px;font-weight:600;cursor:pointer;transition:all 0.2s";
    btn.onmouseover=function(){btn.style.borderColor="#22D3EE";btn.style.boxShadow="0 2px 8px rgba(34,211,238,0.15)"};
    btn.onmouseout=function(){btn.style.borderColor="#E2E8F0";btn.style.boxShadow="none"};
    var self=this;
    btn.onclick=function(){self.authorize()};
    el.appendChild(btn);
  };

  window.AddrPass=AddrPass;
})();`
