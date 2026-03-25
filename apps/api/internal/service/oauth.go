package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/addrpass/addrpass/apps/api/internal/model"
)

var (
	ErrOAuthAppNotFound    = errors.New("oauth app not found")
	ErrInvalidRedirectURI  = errors.New("invalid redirect_uri")
	ErrAuthCodeNotFound    = errors.New("authorization code not found")
	ErrAuthCodeExpired     = errors.New("authorization code expired")
	ErrAuthCodeUsed        = errors.New("authorization code already used")
	ErrClientMismatch      = errors.New("client_id does not match")
	ErrRedirectURIMismatch = errors.New("redirect_uri does not match")
)

type OAuthService struct {
	db        *pgxpool.Pool
	jwtSecret string
}

func NewOAuthService(db *pgxpool.Pool, jwtSecret string) *OAuthService {
	return &OAuthService{db: db, jwtSecret: jwtSecret}
}

// CreateApp registers an OAuth application for a business.
func (s *OAuthService) CreateApp(ctx context.Context, userID, businessID string, req model.CreateOAuthAppRequest) (*model.CreateOAuthAppResponse, error) {
	// Verify business ownership
	var count int
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM businesses WHERE id = $1 AND owner_id = $2`, businessID, userID).Scan(&count)
	if err != nil || count == 0 {
		return nil, ErrBusinessNotFound
	}

	clientID, err := generateClientID()
	if err != nil {
		return nil, err
	}
	clientSecret, err := generateClientSecret()
	if err != nil {
		return nil, err
	}
	secretHash, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var app model.OAuthApp
	err = s.db.QueryRow(ctx,
		`INSERT INTO oauth_apps (business_id, name, logo_url, redirect_uris, client_id, client_secret_hash)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, business_id, name, logo_url, redirect_uris, client_id, active, created_at`,
		businessID, req.Name, req.LogoURL, req.RedirectURIs, clientID, string(secretHash),
	).Scan(&app.ID, &app.BusinessID, &app.Name, &app.LogoURL, &app.RedirectURIs, &app.ClientID, &app.Active, &app.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &model.CreateOAuthAppResponse{App: app, ClientSecret: clientSecret}, nil
}

// GetAppByClientID returns the OAuth app for display on the consent screen.
func (s *OAuthService) GetAppByClientID(ctx context.Context, clientID string) (*model.OAuthApp, error) {
	var app model.OAuthApp
	err := s.db.QueryRow(ctx,
		`SELECT id, business_id, name, logo_url, redirect_uris, client_id, active, created_at
		 FROM oauth_apps WHERE client_id = $1 AND active = TRUE`,
		clientID,
	).Scan(&app.ID, &app.BusinessID, &app.Name, &app.LogoURL, &app.RedirectURIs, &app.ClientID, &app.Active, &app.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOAuthAppNotFound
		}
		return nil, err
	}
	return &app, nil
}

// ValidateRedirectURI checks that the redirect URI is registered for the app.
func (s *OAuthService) ValidateRedirectURI(app *model.OAuthApp, redirectURI string) error {
	for _, uri := range app.RedirectURIs {
		if uri == redirectURI {
			return nil
		}
	}
	return ErrInvalidRedirectURI
}

// CreateAuthorizationCode creates a short-lived authorization code after user consent.
func (s *OAuthService) CreateAuthorizationCode(ctx context.Context, userID string, req model.ConsentRequest) (string, error) {
	app, err := s.GetAppByClientID(ctx, req.ClientID)
	if err != nil {
		return "", err
	}
	if err := s.ValidateRedirectURI(app, req.RedirectURI); err != nil {
		return "", err
	}

	code, err := generateAuthCode()
	if err != nil {
		return "", err
	}

	scope := req.Scope
	if scope == "" {
		scope = "full"
	}

	_, err = s.db.Exec(ctx,
		`INSERT INTO authorization_codes (code, business_id, user_id, address_id, scope, redirect_uri, state, expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		code, app.BusinessID, userID, req.AddressID, scope, req.RedirectURI, req.State,
		time.Now().Add(10*time.Minute),
	)
	if err != nil {
		return "", err
	}

	return code, nil
}

// ExchangeCode exchanges an authorization code for an access token and creates a share.
func (s *OAuthService) ExchangeCode(ctx context.Context, req model.TokenExchangeRequest) (*model.TokenExchangeResponse, error) {
	// Look up the authorization code
	var ac model.AuthorizationCode
	var businessName string
	err := s.db.QueryRow(ctx,
		`SELECT ac.id, ac.code, ac.business_id, ac.user_id, ac.address_id, ac.scope, ac.redirect_uri, ac.state, ac.expires_at, ac.used, b.name
		 FROM authorization_codes ac JOIN businesses b ON ac.business_id = b.id
		 WHERE ac.code = $1`,
		req.Code,
	).Scan(&ac.ID, &ac.Code, &ac.BusinessID, &ac.UserID, &ac.AddressID, &ac.Scope, &ac.RedirectURI, &ac.State, &ac.ExpiresAt, &ac.Used, &businessName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAuthCodeNotFound
		}
		return nil, err
	}

	if ac.Used {
		return nil, ErrAuthCodeUsed
	}
	if time.Now().After(ac.ExpiresAt) {
		return nil, ErrAuthCodeExpired
	}

	// Validate the OAuth app
	var secretHash string
	var appClientID string
	err = s.db.QueryRow(ctx,
		`SELECT client_id, client_secret_hash FROM oauth_apps WHERE business_id = $1 AND active = TRUE LIMIT 1`,
		ac.BusinessID,
	).Scan(&appClientID, &secretHash)
	if err != nil {
		return nil, ErrOAuthAppNotFound
	}

	if req.ClientID != appClientID {
		return nil, ErrClientMismatch
	}
	if err := bcrypt.CompareHashAndPassword([]byte(secretHash), []byte(req.ClientSecret)); err != nil {
		return nil, ErrInvalidClient
	}
	if req.RedirectURI != ac.RedirectURI {
		return nil, ErrRedirectURIMismatch
	}

	// Mark code as used
	_, _ = s.db.Exec(ctx, `UPDATE authorization_codes SET used = TRUE WHERE id = $1`, ac.ID)

	// Create a share on behalf of the user
	shareToken, err := generateToken()
	if err != nil {
		return nil, err
	}

	var shareID string
	err = s.db.QueryRow(ctx,
		`INSERT INTO shares (address_id, user_id, token, access_type, scope, delegated_to_business_id)
		 VALUES ($1, $2, $3, 'authenticated', $4, $5)
		 RETURNING id`,
		ac.AddressID, ac.UserID, shareToken, ac.Scope, ac.BusinessID,
	).Scan(&shareID)
	if err != nil {
		return nil, err
	}

	// Generate a business access token
	expiresIn := 86400 // 24 hours
	claims := jwt.MapClaims{
		"sub":           ac.BusinessID,
		"type":          "business",
		"business_name": businessName,
		"scopes":        []string{string(ac.Scope)},
		"share_id":      shareID,
		"iat":           time.Now().Unix(),
		"exp":           time.Now().Add(time.Duration(expiresIn) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &model.TokenExchangeResponse{
		AccessToken: signed,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		Scope:       string(ac.Scope),
		ShareToken:  shareToken,
	}, nil
}

func generateAuthCode() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("apc_%s", base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)), nil
}
