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

type Share struct {
	ID          string      `json:"id"`
	AddressID   string      `json:"address_id"`
	UserID      string      `json:"user_id"`
	Token       string      `json:"token"`
	AccessType  ShareAccess `json:"access_type"`
	Pin         string      `json:"pin,omitempty"`
	ExpiresAt   *time.Time  `json:"expires_at,omitempty"`
	MaxAccesses *int        `json:"max_accesses,omitempty"`
	AccessCount int         `json:"access_count"`
	Active      bool        `json:"active"`
	CreatedAt   time.Time   `json:"created_at"`
}

type AccessLog struct {
	ID        string    `json:"id"`
	ShareID   string    `json:"share_id"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Country   string    `json:"country,omitempty"`
	AccessAt  time.Time `json:"access_at"`
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
	Pin         string      `json:"pin,omitempty"`
	ExpiresAt   *time.Time  `json:"expires_at,omitempty"`
	MaxAccesses *int        `json:"max_accesses,omitempty"`
}

type ResolveResponse struct {
	Address Address `json:"address"`
}
