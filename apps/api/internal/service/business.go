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
	ErrBusinessNotFound = errors.New("business not found")
	ErrAPIKeyNotFound   = errors.New("api key not found")
	ErrInvalidClient    = errors.New("invalid client credentials")
	ErrInvalidGrant     = errors.New("unsupported grant type")
)

type BusinessService struct {
	db        *pgxpool.Pool
	jwtSecret string
}

func NewBusinessService(db *pgxpool.Pool, jwtSecret string) *BusinessService {
	return &BusinessService{db: db, jwtSecret: jwtSecret}
}

func (s *BusinessService) CreateBusiness(ctx context.Context, userID string, req model.CreateBusinessRequest) (*model.Business, error) {
	var biz model.Business
	err := s.db.QueryRow(ctx,
		`INSERT INTO businesses (name, owner_id) VALUES ($1, $2)
		 RETURNING id, name, owner_id, created_at, updated_at`,
		req.Name, userID,
	).Scan(&biz.ID, &biz.Name, &biz.OwnerID, &biz.CreatedAt, &biz.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &biz, nil
}

func (s *BusinessService) ListBusinesses(ctx context.Context, userID string) ([]model.Business, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, name, owner_id, created_at, updated_at FROM businesses WHERE owner_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var businesses []model.Business
	for rows.Next() {
		var b model.Business
		if err := rows.Scan(&b.ID, &b.Name, &b.OwnerID, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		businesses = append(businesses, b)
	}
	if businesses == nil {
		businesses = []model.Business{}
	}
	return businesses, nil
}

func (s *BusinessService) CreateAPIKey(ctx context.Context, userID, businessID string, req model.CreateAPIKeyRequest) (*model.CreateAPIKeyResponse, error) {
	// Verify ownership
	var count int
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM businesses WHERE id = $1 AND owner_id = $2`, businessID, userID).Scan(&count)
	if err != nil || count == 0 {
		return nil, ErrBusinessNotFound
	}

	// Generate client_id and client_secret
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

	scopes := req.Scopes
	if len(scopes) == 0 {
		scopes = []string{"full"}
	}

	var key model.APIKey
	err = s.db.QueryRow(ctx,
		`INSERT INTO api_keys (business_id, client_id, client_secret_hash, name, scopes)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, business_id, client_id, name, scopes, rate_limit_per_hour, active, created_at`,
		businessID, clientID, string(secretHash), req.Name, scopes,
	).Scan(&key.ID, &key.BusinessID, &key.ClientID, &key.Name, &key.Scopes, &key.RateLimitPerHour, &key.Active, &key.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &model.CreateAPIKeyResponse{
		APIKey:       key,
		ClientSecret: clientSecret,
	}, nil
}

func (s *BusinessService) ListAPIKeys(ctx context.Context, userID, businessID string) ([]model.APIKey, error) {
	// Verify ownership
	var count int
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM businesses WHERE id = $1 AND owner_id = $2`, businessID, userID).Scan(&count)
	if err != nil || count == 0 {
		return nil, ErrBusinessNotFound
	}

	rows, err := s.db.Query(ctx,
		`SELECT id, business_id, client_id, name, scopes, rate_limit_per_hour, active, last_used_at, created_at
		 FROM api_keys WHERE business_id = $1 ORDER BY created_at DESC`,
		businessID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []model.APIKey
	for rows.Next() {
		var k model.APIKey
		if err := rows.Scan(&k.ID, &k.BusinessID, &k.ClientID, &k.Name, &k.Scopes, &k.RateLimitPerHour, &k.Active, &k.LastUsedAt, &k.CreatedAt); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	if keys == nil {
		keys = []model.APIKey{}
	}
	return keys, nil
}

func (s *BusinessService) RevokeAPIKey(ctx context.Context, userID, keyID string) error {
	result, err := s.db.Exec(ctx,
		`UPDATE api_keys SET active = FALSE
		 WHERE id = $1 AND business_id IN (SELECT id FROM businesses WHERE owner_id = $2)`,
		keyID, userID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrAPIKeyNotFound
	}
	return nil
}

// OAuthToken handles the client_credentials grant flow.
// Returns a JWT scoped to the business with the API key's scopes.
func (s *BusinessService) OAuthToken(ctx context.Context, req model.OAuthTokenRequest) (*model.OAuthTokenResponse, error) {
	if req.GrantType != "client_credentials" {
		return nil, ErrInvalidGrant
	}

	// Lookup API key by client_id
	var keyID, businessID, secretHash, businessName string
	var scopes []string
	var active bool
	err := s.db.QueryRow(ctx,
		`SELECT ak.id, ak.business_id, ak.client_secret_hash, ak.scopes, ak.active, b.name
		 FROM api_keys ak JOIN businesses b ON ak.business_id = b.id
		 WHERE ak.client_id = $1`,
		req.ClientID,
	).Scan(&keyID, &businessID, &secretHash, &scopes, &active, &businessName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidClient
		}
		return nil, err
	}

	if !active {
		return nil, ErrInvalidClient
	}

	if err := bcrypt.CompareHashAndPassword([]byte(secretHash), []byte(req.ClientSecret)); err != nil {
		return nil, ErrInvalidClient
	}

	// Update last_used_at
	_, _ = s.db.Exec(ctx, `UPDATE api_keys SET last_used_at = NOW() WHERE id = $1`, keyID)

	// Generate JWT with business claims
	expiresIn := 3600 // 1 hour
	claims := jwt.MapClaims{
		"sub":           businessID,
		"type":          "business",
		"business_name": businessName,
		"scopes":        scopes,
		"key_id":        keyID,
		"iat":           time.Now().Unix(),
		"exp":           time.Now().Add(time.Duration(expiresIn) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	scopeStr := ""
	for i, sc := range scopes {
		if i > 0 {
			scopeStr += " "
		}
		scopeStr += sc
	}

	return &model.OAuthTokenResponse{
		AccessToken: signed,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		Scope:       scopeStr,
	}, nil
}

// ValidateBusinessToken parses a business JWT and returns business info.
func (s *BusinessService) ValidateBusinessToken(tokenStr, jwtSecret string) (businessID, businessName string, scopes []string, err error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", "", nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", nil, fmt.Errorf("invalid claims")
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != "business" {
		return "", "", nil, fmt.Errorf("not a business token")
	}

	businessID, _ = claims["sub"].(string)
	businessName, _ = claims["business_name"].(string)

	if scopesRaw, ok := claims["scopes"].([]interface{}); ok {
		for _, s := range scopesRaw {
			if str, ok := s.(string); ok {
				scopes = append(scopes, str)
			}
		}
	}

	return businessID, businessName, scopes, nil
}

func generateClientID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "ap_" + base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}

func generateClientSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "aps_" + base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}
