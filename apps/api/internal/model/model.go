package model

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Address struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Label     string    `json:"label"`
	Line1     string    `json:"line1"`
	Line2     string    `json:"line2,omitempty"`
	City      string    `json:"city"`
	State     string    `json:"state,omitempty"`
	PostCode  string    `json:"post_code"`
	Country   string    `json:"country"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ShareAccess string

const (
	ShareAccessPublic        ShareAccess = "public"
	ShareAccessAuthenticated ShareAccess = "authenticated"
)

type ShareScope string

const (
	ScopeFull     ShareScope = "full"
	ScopeDelivery ShareScope = "delivery"
	ScopeZone     ShareScope = "zone"
	ScopeVerify   ShareScope = "verify"
)

type Share struct {
	ID          string      `json:"id"`
	AddressID   string      `json:"address_id"`
	UserID      string      `json:"user_id"`
	Token       string      `json:"token"`
	AccessType  ShareAccess `json:"access_type"`
	Scope       ShareScope  `json:"scope"`
	Pin         string      `json:"pin,omitempty"`
	ExpiresAt   *time.Time  `json:"expires_at,omitempty"`
	MaxAccesses *int        `json:"max_accesses,omitempty"`
	AccessCount int         `json:"access_count"`
	Active      bool        `json:"active"`
	CreatedAt   time.Time   `json:"created_at"`
}

type AccessLog struct {
	ID           string    `json:"id"`
	ShareID      string    `json:"share_id"`
	IP           string    `json:"ip"`
	UserAgent    string    `json:"user_agent"`
	Country      string    `json:"country,omitempty"`
	BusinessID   string    `json:"business_id,omitempty"`
	BusinessName string    `json:"business_name,omitempty"`
	AccessAt     time.Time `json:"access_at"`
}

// Business & API key types

type Business struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type APIKey struct {
	ID               string    `json:"id"`
	BusinessID       string    `json:"business_id"`
	ClientID         string    `json:"client_id"`
	Name             string    `json:"name"`
	Scopes           []string  `json:"scopes"`
	RateLimitPerHour int       `json:"rate_limit_per_hour"`
	Active           bool      `json:"active"`
	LastUsedAt       *time.Time `json:"last_used_at,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type Webhook struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	URL             string     `json:"url"`
	Secret          string     `json:"secret,omitempty"`
	Events          []string   `json:"events"`
	Active          bool       `json:"active"`
	FailureCount    int        `json:"failure_count"`
	LastTriggeredAt *time.Time `json:"last_triggered_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// Request/Response types

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateAddressRequest struct {
	Label    string `json:"label"`
	Line1    string `json:"line1"`
	Line2    string `json:"line2,omitempty"`
	City     string `json:"city"`
	State    string `json:"state,omitempty"`
	PostCode string `json:"post_code"`
	Country  string `json:"country"`
	Phone    string `json:"phone,omitempty"`
}

type UpdateAddressRequest struct {
	Label    *string `json:"label,omitempty"`
	Line1    *string `json:"line1,omitempty"`
	Line2    *string `json:"line2,omitempty"`
	City     *string `json:"city,omitempty"`
	State    *string `json:"state,omitempty"`
	PostCode *string `json:"post_code,omitempty"`
	Country  *string `json:"country,omitempty"`
	Phone    *string `json:"phone,omitempty"`
}

type CreateShareRequest struct {
	AddressID   string      `json:"address_id"`
	AccessType  ShareAccess `json:"access_type"`
	Scope       ShareScope  `json:"scope,omitempty"`
	Pin         string      `json:"pin,omitempty"`
	ExpiresAt   *time.Time  `json:"expires_at,omitempty"`
	MaxAccesses *int        `json:"max_accesses,omitempty"`
}

type ResolveResponse struct {
	Address Address `json:"address"`
	Scope   string  `json:"scope"`
}

// Business requests

type CreateBusinessRequest struct {
	Name string `json:"name"`
}

type CreateAPIKeyRequest struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes,omitempty"`
}

type CreateAPIKeyResponse struct {
	APIKey       APIKey `json:"api_key"`
	ClientSecret string `json:"client_secret"`
}

type OAuthTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// Webhook requests

type CreateWebhookRequest struct {
	URL    string   `json:"url"`
	Events []string `json:"events,omitempty"`
}

type WebhookEvent struct {
	Event     string      `json:"event"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type AccessEvent struct {
	ShareID      string `json:"share_id"`
	Token        string `json:"token"`
	IP           string `json:"ip"`
	UserAgent    string `json:"user_agent"`
	BusinessName string `json:"business_name,omitempty"`
	Scope        string `json:"scope"`
}
