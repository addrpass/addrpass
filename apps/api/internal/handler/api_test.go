package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/addrpass/addrpass/apps/api/internal/database"
	"github.com/addrpass/addrpass/apps/api/internal/handler"
	"github.com/addrpass/addrpass/apps/api/internal/middleware"
	"github.com/addrpass/addrpass/apps/api/internal/model"
	"github.com/addrpass/addrpass/apps/api/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

const testJWTSecret = "test-secret-key"
const testBaseURL = "http://localhost:8080"

type testEnv struct {
	pool    *pgxpool.Pool
	router  chi.Router
	authSvc *service.AuthService
}

func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://addrpass@localhost:5432/addrpass?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := database.Connect(ctx, dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean tables before each test
	pool.Exec(ctx, "DELETE FROM authorization_codes")
	pool.Exec(ctx, "DELETE FROM oauth_apps")
	pool.Exec(ctx, "DELETE FROM labels")
	pool.Exec(ctx, "DELETE FROM delegations")
	pool.Exec(ctx, "DELETE FROM webhook_deliveries")
	pool.Exec(ctx, "DELETE FROM webhooks")
	pool.Exec(ctx, "DELETE FROM access_logs")
	pool.Exec(ctx, "DELETE FROM shares")
	pool.Exec(ctx, "DELETE FROM api_keys")
	pool.Exec(ctx, "DELETE FROM businesses")
	pool.Exec(ctx, "DELETE FROM addresses")
	pool.Exec(ctx, "DELETE FROM users")

	authSvc := service.NewAuthService(pool, testJWTSecret)
	addressSvc := service.NewAddressService(pool)
	shareSvc := service.NewShareService(pool)
	webhookSvc := service.NewWebhookService(pool)
	businessSvc := service.NewBusinessService(pool, testJWTSecret)
	delegationSvc := service.NewDelegationService(pool)
	labelSvc := service.NewLabelService(pool)
	oauthSvc := service.NewOAuthService(pool, testJWTSecret)

	authH := handler.NewAuthHandler(authSvc)
	addressH := handler.NewAddressHandler(addressSvc)
	billingSvc := service.NewBillingService(pool)
	shareH := handler.NewShareHandler(shareSvc, webhookSvc, billingSvc, testBaseURL, testJWTSecret)
	businessH := handler.NewBusinessHandler(businessSvc)
	webhookH := handler.NewWebhookHandler(webhookSvc)
	delegationH := handler.NewDelegationHandler(delegationSvc)
	labelH := handler.NewLabelHandler(labelSvc)
	oauthH := handler.NewOAuthHandler(oauthSvc, addressSvc)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"*"}}))

	r.Post("/api/v1/auth/register", authH.Register)
	r.Post("/api/v1/auth/login", authH.Login)
	r.Get("/api/v1/resolve/{token}", shareH.Resolve)
	r.Get("/api/v1/qr/{token}", shareH.QRCode)

	r.Post("/api/v1/oauth/token", businessH.OAuthToken)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(testJWTSecret))
		r.Get("/api/v1/auth/me", authH.Me)
		r.Post("/api/v1/addresses", addressH.Create)
		r.Get("/api/v1/addresses", addressH.List)
		r.Get("/api/v1/addresses/{id}", addressH.Get)
		r.Put("/api/v1/addresses/{id}", addressH.Update)
		r.Delete("/api/v1/addresses/{id}", addressH.Delete)
		r.Post("/api/v1/shares", shareH.Create)
		r.Get("/api/v1/shares", shareH.List)
		r.Patch("/api/v1/shares/{id}/revoke", shareH.Revoke)
		r.Delete("/api/v1/shares/{id}", shareH.Delete)
		r.Get("/api/v1/shares/{id}/accesses", shareH.AccessLogs)
		r.Post("/api/v1/businesses", businessH.CreateBusiness)
		r.Get("/api/v1/businesses", businessH.ListBusinesses)
		r.Post("/api/v1/businesses/{businessId}/api-keys", businessH.CreateAPIKey)
		r.Get("/api/v1/businesses/{businessId}/api-keys", businessH.ListAPIKeys)
		r.Patch("/api/v1/api-keys/{keyId}/revoke", businessH.RevokeAPIKey)
		r.Post("/api/v1/webhooks", webhookH.Create)
		r.Get("/api/v1/webhooks", webhookH.List)
		r.Delete("/api/v1/webhooks/{id}", webhookH.Delete)
		r.Post("/api/v1/delegations", delegationH.CreateByUser)
		r.Get("/api/v1/shares/{shareId}/delegations", delegationH.ListForShare)
		r.Patch("/api/v1/delegations/{id}/revoke", delegationH.Revoke)
		r.Post("/api/v1/labels", labelH.Create)
		r.Post("/api/v1/businesses/{businessId}/oauth-apps", oauthH.CreateApp)
		r.Get("/api/v1/oauth/authorize", oauthH.Authorize)
		r.Post("/api/v1/oauth/consent", oauthH.Consent)
	})
	r.Get("/api/v1/labels/{ref}/image", labelH.GetLabelImage)
	r.Post("/api/v1/oauth/exchange", oauthH.Exchange)

	t.Cleanup(func() { pool.Close() })

	return &testEnv{pool: pool, router: r, authSvc: authSvc}
}

func (e *testEnv) request(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	e.router.ServeHTTP(w, req)
	return w
}

func (e *testEnv) registerUser(t *testing.T, email, password, name string) (string, model.User) {
	t.Helper()
	w := e.request("POST", "/api/v1/auth/register", model.RegisterRequest{
		Email: email, Password: password, Name: name,
	}, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("register failed: %d %s", w.Code, w.Body.String())
	}
	var resp model.AuthResponse
	json.NewDecoder(w.Body).Decode(&resp)
	return resp.Token, resp.User
}

// ─── Auth Tests ─────────────────────────────────────────────

func TestRegister(t *testing.T) {
	env := setupTestEnv(t)

	w := env.request("POST", "/api/v1/auth/register", model.RegisterRequest{
		Email: "test@example.com", Password: "password123", Name: "Test",
	}, "")

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.AuthResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Token == "" {
		t.Fatal("expected token")
	}
	if resp.User.Email != "test@example.com" {
		t.Fatalf("expected email test@example.com, got %s", resp.User.Email)
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	env := setupTestEnv(t)

	env.registerUser(t, "dup@example.com", "password123", "First")

	w := env.request("POST", "/api/v1/auth/register", model.RegisterRequest{
		Email: "dup@example.com", Password: "password123", Name: "Second",
	}, "")

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestRegisterWeakPassword(t *testing.T) {
	env := setupTestEnv(t)

	w := env.request("POST", "/api/v1/auth/register", model.RegisterRequest{
		Email: "weak@example.com", Password: "short", Name: "Weak",
	}, "")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestLogin(t *testing.T) {
	env := setupTestEnv(t)
	env.registerUser(t, "login@example.com", "password123", "Login Test")

	w := env.request("POST", "/api/v1/auth/login", model.LoginRequest{
		Email: "login@example.com", Password: "password123",
	}, "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.AuthResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Token == "" {
		t.Fatal("expected token")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	env := setupTestEnv(t)
	env.registerUser(t, "wrong@example.com", "password123", "Wrong")

	w := env.request("POST", "/api/v1/auth/login", model.LoginRequest{
		Email: "wrong@example.com", Password: "wrongpassword",
	}, "")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMe(t *testing.T) {
	env := setupTestEnv(t)
	token, user := env.registerUser(t, "me@example.com", "password123", "Me Test")

	w := env.request("GET", "/api/v1/auth/me", nil, token)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var got model.User
	json.NewDecoder(w.Body).Decode(&got)
	if got.ID != user.ID {
		t.Fatalf("expected user ID %s, got %s", user.ID, got.ID)
	}
}

func TestMeUnauthorized(t *testing.T) {
	env := setupTestEnv(t)

	w := env.request("GET", "/api/v1/auth/me", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// ─── Address Tests ──────────────────────────────────────────

func TestAddressCRUD(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := env.registerUser(t, "addr@example.com", "password123", "Addr Test")

	// Create
	w := env.request("POST", "/api/v1/addresses", model.CreateAddressRequest{
		Label: "Home", Line1: "123 Main St", City: "Istanbul", PostCode: "34000", Country: "TR",
	}, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var addr model.Address
	json.NewDecoder(w.Body).Decode(&addr)
	if addr.Label != "Home" {
		t.Fatalf("expected label Home, got %s", addr.Label)
	}

	// List
	w = env.request("GET", "/api/v1/addresses", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", w.Code)
	}
	var addresses []model.Address
	json.NewDecoder(w.Body).Decode(&addresses)
	if len(addresses) != 1 {
		t.Fatalf("expected 1 address, got %d", len(addresses))
	}

	// Get
	w = env.request("GET", "/api/v1/addresses/"+addr.ID, nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", w.Code)
	}

	// Update
	newLabel := "Work"
	w = env.request("PUT", "/api/v1/addresses/"+addr.ID, model.UpdateAddressRequest{
		Label: &newLabel,
	}, token)
	if w.Code != http.StatusOK {
		t.Fatalf("update: expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var updated model.Address
	json.NewDecoder(w.Body).Decode(&updated)
	if updated.Label != "Work" {
		t.Fatalf("expected label Work, got %s", updated.Label)
	}

	// Delete
	w = env.request("DELETE", "/api/v1/addresses/"+addr.ID, nil, token)
	if w.Code != http.StatusNoContent {
		t.Fatalf("delete: expected 204, got %d", w.Code)
	}

	// Verify deleted
	w = env.request("GET", "/api/v1/addresses/"+addr.ID, nil, token)
	if w.Code != http.StatusNotFound {
		t.Fatalf("after delete: expected 404, got %d", w.Code)
	}
}

func TestAddressIsolation(t *testing.T) {
	env := setupTestEnv(t)
	token1, _ := env.registerUser(t, "user1@example.com", "password123", "User1")
	token2, _ := env.registerUser(t, "user2@example.com", "password123", "User2")

	// User1 creates an address
	w := env.request("POST", "/api/v1/addresses", model.CreateAddressRequest{
		Label: "Secret", Line1: "Hidden St", City: "Private", PostCode: "00000", Country: "XX",
	}, token1)
	var addr model.Address
	json.NewDecoder(w.Body).Decode(&addr)

	// User2 cannot access it
	w = env.request("GET", "/api/v1/addresses/"+addr.ID, nil, token2)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d — user isolation broken", w.Code)
	}

	// User2 cannot delete it
	w = env.request("DELETE", "/api/v1/addresses/"+addr.ID, nil, token2)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d — user isolation broken", w.Code)
	}
}

// ─── Share Tests ────────────────────────────────────────────

func TestShareFlow(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := env.registerUser(t, "share@example.com", "password123", "Share Test")

	// Create address
	w := env.request("POST", "/api/v1/addresses", model.CreateAddressRequest{
		Label: "Home", Line1: "456 Token Ave", City: "Berlin", PostCode: "10115", Country: "DE",
	}, token)
	var addr model.Address
	json.NewDecoder(w.Body).Decode(&addr)

	// Create share
	w = env.request("POST", "/api/v1/shares", model.CreateShareRequest{
		AddressID: addr.ID, AccessType: model.ShareAccessPublic,
	}, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create share: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var shareResp struct {
		Share model.Share `json:"share"`
		URL   string      `json:"url"`
	}
	json.NewDecoder(w.Body).Decode(&shareResp)
	if shareResp.Share.Token == "" {
		t.Fatal("expected share token")
	}
	if shareResp.URL == "" {
		t.Fatal("expected share URL")
	}

	// Resolve token — no auth needed
	w = env.request("GET", "/api/v1/resolve/"+shareResp.Share.Token, nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("resolve: expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resolveResp model.ResolveResponse
	json.NewDecoder(w.Body).Decode(&resolveResp)
	if resolveResp.Address.Line1 != "456 Token Ave" {
		t.Fatalf("expected '456 Token Ave', got '%s'", resolveResp.Address.Line1)
	}

	// Check access logs
	w = env.request("GET", "/api/v1/shares/"+shareResp.Share.ID+"/accesses", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("access logs: expected 200, got %d", w.Code)
	}
	var logs []model.AccessLog
	json.NewDecoder(w.Body).Decode(&logs)
	if len(logs) != 1 {
		t.Fatalf("expected 1 access log, got %d", len(logs))
	}

	// List shares
	w = env.request("GET", "/api/v1/shares", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("list shares: expected 200, got %d", w.Code)
	}
	var shares []model.Share
	json.NewDecoder(w.Body).Decode(&shares)
	if len(shares) != 1 {
		t.Fatalf("expected 1 share, got %d", len(shares))
	}

	// Revoke
	w = env.request("PATCH", "/api/v1/shares/"+shareResp.Share.ID+"/revoke", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("revoke: expected 200, got %d", w.Code)
	}

	// Resolve after revoke — should fail
	w = env.request("GET", "/api/v1/resolve/"+shareResp.Share.Token, nil, "")
	if w.Code != http.StatusGone {
		t.Fatalf("resolve after revoke: expected 410, got %d", w.Code)
	}
}

func TestShareWithPin(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := env.registerUser(t, "pin@example.com", "password123", "Pin Test")

	w := env.request("POST", "/api/v1/addresses", model.CreateAddressRequest{
		Label: "Home", Line1: "789 Pin St", City: "Paris", PostCode: "75001", Country: "FR",
	}, token)
	var addr model.Address
	json.NewDecoder(w.Body).Decode(&addr)

	// Create share with PIN
	w = env.request("POST", "/api/v1/shares", model.CreateShareRequest{
		AddressID: addr.ID, AccessType: model.ShareAccessPublic, Pin: "1234",
	}, token)
	var shareResp struct {
		Share model.Share `json:"share"`
	}
	json.NewDecoder(w.Body).Decode(&shareResp)

	// Resolve without PIN — should fail
	w = env.request("GET", "/api/v1/resolve/"+shareResp.Share.Token, nil, "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("no pin: expected 403, got %d", w.Code)
	}

	// Resolve with wrong PIN
	w = env.request("GET", "/api/v1/resolve/"+shareResp.Share.Token+"?pin=0000", nil, "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("wrong pin: expected 403, got %d", w.Code)
	}

	// Resolve with correct PIN
	w = env.request("GET", "/api/v1/resolve/"+shareResp.Share.Token+"?pin=1234", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("correct pin: expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestShareMaxAccesses(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := env.registerUser(t, "max@example.com", "password123", "Max Test")

	w := env.request("POST", "/api/v1/addresses", model.CreateAddressRequest{
		Label: "Home", Line1: "1 Max St", City: "London", PostCode: "EC1A", Country: "GB",
	}, token)
	var addr model.Address
	json.NewDecoder(w.Body).Decode(&addr)

	// Create share with max 2 accesses
	maxAccesses := 2
	w = env.request("POST", "/api/v1/shares", model.CreateShareRequest{
		AddressID: addr.ID, AccessType: model.ShareAccessPublic, MaxAccesses: &maxAccesses,
	}, token)
	var shareResp struct {
		Share model.Share `json:"share"`
	}
	json.NewDecoder(w.Body).Decode(&shareResp)

	// Access 1 — OK
	w = env.request("GET", "/api/v1/resolve/"+shareResp.Share.Token, nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("access 1: expected 200, got %d", w.Code)
	}

	// Access 2 — OK
	w = env.request("GET", "/api/v1/resolve/"+shareResp.Share.Token, nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("access 2: expected 200, got %d", w.Code)
	}

	// Access 3 — should fail
	w = env.request("GET", "/api/v1/resolve/"+shareResp.Share.Token, nil, "")
	if w.Code != http.StatusGone {
		t.Fatalf("access 3: expected 410, got %d", w.Code)
	}
}

func TestQRCode(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := env.registerUser(t, "qr@example.com", "password123", "QR Test")

	w := env.request("POST", "/api/v1/addresses", model.CreateAddressRequest{
		Label: "Home", Line1: "QR Lane", City: "Tokyo", PostCode: "100-0001", Country: "JP",
	}, token)
	var addr model.Address
	json.NewDecoder(w.Body).Decode(&addr)

	w = env.request("POST", "/api/v1/shares", model.CreateShareRequest{
		AddressID: addr.ID, AccessType: model.ShareAccessPublic,
	}, token)
	var shareResp struct {
		Share model.Share `json:"share"`
	}
	json.NewDecoder(w.Body).Decode(&shareResp)

	// Get QR code
	w = env.request("GET", "/api/v1/qr/"+shareResp.Share.Token, nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("qr: expected 200, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "image/png" {
		t.Fatalf("expected image/png, got %s", w.Header().Get("Content-Type"))
	}
	if w.Body.Len() < 100 {
		t.Fatal("QR code PNG seems too small")
	}
}

func TestResolveNotFound(t *testing.T) {
	env := setupTestEnv(t)

	w := env.request("GET", "/api/v1/resolve/nonexistent-token", nil, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// ─── Phase 2: Business & API Key Tests ──────────────────────

func TestBusinessAndAPIKeyFlow(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := env.registerUser(t, "biz@example.com", "password123", "Biz Owner")

	// Create business
	w := env.request("POST", "/api/v1/businesses", model.CreateBusinessRequest{
		Name: "DHL Test",
	}, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create business: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var biz model.Business
	json.NewDecoder(w.Body).Decode(&biz)
	if biz.Name != "DHL Test" {
		t.Fatalf("expected name 'DHL Test', got '%s'", biz.Name)
	}

	// List businesses
	w = env.request("GET", "/api/v1/businesses", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("list businesses: expected 200, got %d", w.Code)
	}
	var businesses []model.Business
	json.NewDecoder(w.Body).Decode(&businesses)
	if len(businesses) != 1 {
		t.Fatalf("expected 1 business, got %d", len(businesses))
	}

	// Create API key
	w = env.request("POST", "/api/v1/businesses/"+biz.ID+"/api-keys", model.CreateAPIKeyRequest{
		Name:   "Production Key",
		Scopes: []string{"delivery"},
	}, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create api key: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var keyResp model.CreateAPIKeyResponse
	json.NewDecoder(w.Body).Decode(&keyResp)
	if keyResp.ClientSecret == "" {
		t.Fatal("expected client_secret")
	}
	if keyResp.APIKey.ClientID == "" {
		t.Fatal("expected client_id")
	}

	// List API keys
	w = env.request("GET", "/api/v1/businesses/"+biz.ID+"/api-keys", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("list api keys: expected 200, got %d", w.Code)
	}
	var keys []model.APIKey
	json.NewDecoder(w.Body).Decode(&keys)
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}

	// OAuth token exchange
	w = env.request("POST", "/api/v1/oauth/token", model.OAuthTokenRequest{
		GrantType:    "client_credentials",
		ClientID:     keyResp.APIKey.ClientID,
		ClientSecret: keyResp.ClientSecret,
	}, "")
	if w.Code != http.StatusOK {
		t.Fatalf("oauth token: expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var oauthResp model.OAuthTokenResponse
	json.NewDecoder(w.Body).Decode(&oauthResp)
	if oauthResp.AccessToken == "" {
		t.Fatal("expected access_token")
	}
	if oauthResp.TokenType != "Bearer" {
		t.Fatalf("expected Bearer, got %s", oauthResp.TokenType)
	}

	// Revoke API key
	w = env.request("PATCH", "/api/v1/api-keys/"+keyResp.APIKey.ID+"/revoke", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("revoke key: expected 200, got %d", w.Code)
	}

	// OAuth should fail now
	w = env.request("POST", "/api/v1/oauth/token", model.OAuthTokenRequest{
		GrantType:    "client_credentials",
		ClientID:     keyResp.APIKey.ClientID,
		ClientSecret: keyResp.ClientSecret,
	}, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("oauth after revoke: expected 401, got %d", w.Code)
	}
}

func TestOAuthInvalidCredentials(t *testing.T) {
	env := setupTestEnv(t)

	w := env.request("POST", "/api/v1/oauth/token", model.OAuthTokenRequest{
		GrantType:    "client_credentials",
		ClientID:     "nonexistent",
		ClientSecret: "wrong",
	}, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestOAuthInvalidGrantType(t *testing.T) {
	env := setupTestEnv(t)

	w := env.request("POST", "/api/v1/oauth/token", model.OAuthTokenRequest{
		GrantType: "authorization_code",
	}, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// ─── Phase 2: Scoped Access Tests ───────────────────────────

func TestScopedAccess(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := env.registerUser(t, "scope@example.com", "password123", "Scope Test")

	// Create address with all fields
	w := env.request("POST", "/api/v1/addresses", model.CreateAddressRequest{
		Label: "Home", Line1: "100 Scope St", Line2: "Apt 5", City: "Munich",
		State: "Bavaria", PostCode: "80331", Country: "DE", Phone: "+49123456",
	}, token)
	var addr model.Address
	json.NewDecoder(w.Body).Decode(&addr)

	// Test "full" scope — all fields returned
	w = env.request("POST", "/api/v1/shares", model.CreateShareRequest{
		AddressID: addr.ID, AccessType: model.ShareAccessPublic, Scope: model.ScopeFull,
	}, token)
	var fullShare struct{ Share model.Share `json:"share"` }
	json.NewDecoder(w.Body).Decode(&fullShare)

	w = env.request("GET", "/api/v1/resolve/"+fullShare.Share.Token, nil, "")
	var fullResp model.ResolveResponse
	json.NewDecoder(w.Body).Decode(&fullResp)
	if fullResp.Address.Phone != "+49123456" {
		t.Fatalf("full scope: expected phone, got '%s'", fullResp.Address.Phone)
	}
	if fullResp.Scope != "full" {
		t.Fatalf("expected scope 'full', got '%s'", fullResp.Scope)
	}

	// Test "delivery" scope — no phone
	w = env.request("POST", "/api/v1/shares", model.CreateShareRequest{
		AddressID: addr.ID, AccessType: model.ShareAccessPublic, Scope: model.ScopeDelivery,
	}, token)
	var delShare struct{ Share model.Share `json:"share"` }
	json.NewDecoder(w.Body).Decode(&delShare)

	w = env.request("GET", "/api/v1/resolve/"+delShare.Share.Token, nil, "")
	var delResp model.ResolveResponse
	json.NewDecoder(w.Body).Decode(&delResp)
	if delResp.Address.Phone != "" {
		t.Fatalf("delivery scope: phone should be empty, got '%s'", delResp.Address.Phone)
	}
	if delResp.Address.Line1 != "100 Scope St" {
		t.Fatalf("delivery scope: expected line1, got '%s'", delResp.Address.Line1)
	}

	// Test "zone" scope — only city, state, post_code, country
	w = env.request("POST", "/api/v1/shares", model.CreateShareRequest{
		AddressID: addr.ID, AccessType: model.ShareAccessPublic, Scope: model.ScopeZone,
	}, token)
	var zoneShare struct{ Share model.Share `json:"share"` }
	json.NewDecoder(w.Body).Decode(&zoneShare)

	w = env.request("GET", "/api/v1/resolve/"+zoneShare.Share.Token, nil, "")
	var zoneResp model.ResolveResponse
	json.NewDecoder(w.Body).Decode(&zoneResp)
	if zoneResp.Address.Line1 != "" {
		t.Fatalf("zone scope: line1 should be empty, got '%s'", zoneResp.Address.Line1)
	}
	if zoneResp.Address.City != "Munich" {
		t.Fatalf("zone scope: expected city Munich, got '%s'", zoneResp.Address.City)
	}

	// Test "verify" scope — only country
	w = env.request("POST", "/api/v1/shares", model.CreateShareRequest{
		AddressID: addr.ID, AccessType: model.ShareAccessPublic, Scope: model.ScopeVerify,
	}, token)
	var verifyShare struct{ Share model.Share `json:"share"` }
	json.NewDecoder(w.Body).Decode(&verifyShare)

	w = env.request("GET", "/api/v1/resolve/"+verifyShare.Share.Token, nil, "")
	var verifyResp model.ResolveResponse
	json.NewDecoder(w.Body).Decode(&verifyResp)
	if verifyResp.Address.City != "" {
		t.Fatalf("verify scope: city should be empty, got '%s'", verifyResp.Address.City)
	}
	if verifyResp.Address.Country != "DE" {
		t.Fatalf("verify scope: expected country DE, got '%s'", verifyResp.Address.Country)
	}
}

// ─── Phase 2: Webhook Tests ────────────────────────────────

func TestWebhookCRUD(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := env.registerUser(t, "wh@example.com", "password123", "Webhook Test")

	// Create webhook
	w := env.request("POST", "/api/v1/webhooks", model.CreateWebhookRequest{
		URL: "https://example.com/webhook",
	}, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create webhook: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var wh model.Webhook
	json.NewDecoder(w.Body).Decode(&wh)
	if wh.URL != "https://example.com/webhook" {
		t.Fatalf("expected URL, got '%s'", wh.URL)
	}
	if wh.Secret == "" {
		t.Fatal("expected webhook secret")
	}

	// List webhooks
	w = env.request("GET", "/api/v1/webhooks", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("list webhooks: expected 200, got %d", w.Code)
	}
	var webhooks []model.Webhook
	json.NewDecoder(w.Body).Decode(&webhooks)
	if len(webhooks) != 1 {
		t.Fatalf("expected 1 webhook, got %d", len(webhooks))
	}

	// Delete webhook
	w = env.request("DELETE", "/api/v1/webhooks/"+wh.ID, nil, token)
	if w.Code != http.StatusNoContent {
		t.Fatalf("delete webhook: expected 204, got %d", w.Code)
	}

	// Verify deleted
	w = env.request("GET", "/api/v1/webhooks", nil, token)
	json.NewDecoder(w.Body).Decode(&webhooks)
	if len(webhooks) != 0 {
		t.Fatalf("expected 0 webhooks, got %d", len(webhooks))
	}
}

// ─── Phase 3: Delegation Tests ──────────────────────────────

func TestDelegationFlow(t *testing.T) {
	env := setupTestEnv(t)
	userToken, _ := env.registerUser(t, "owner@example.com", "password123", "Owner")

	// Create address + share
	w := env.request("POST", "/api/v1/addresses", model.CreateAddressRequest{
		Label: "Home", Line1: "10 Delegation Rd", City: "Berlin", PostCode: "10115", Country: "DE", Phone: "+4930123",
	}, userToken)
	var addr model.Address
	json.NewDecoder(w.Body).Decode(&addr)

	w = env.request("POST", "/api/v1/shares", model.CreateShareRequest{
		AddressID: addr.ID, AccessType: model.ShareAccessPublic, Scope: model.ScopeFull,
	}, userToken)
	var shareResp struct{ Share model.Share `json:"share"` }
	json.NewDecoder(w.Body).Decode(&shareResp)

	// Create a business to delegate to
	w = env.request("POST", "/api/v1/businesses", model.CreateBusinessRequest{Name: "DHL Express"}, userToken)
	var biz model.Business
	json.NewDecoder(w.Body).Decode(&biz)

	// Delegate share to business with "delivery" scope
	w = env.request("POST", "/api/v1/delegations", model.CreateDelegationRequest{
		ShareID: shareResp.Share.ID, ToBusinessID: biz.ID, Scope: model.ScopeDelivery, Note: "For package #1234",
	}, userToken)
	if w.Code != http.StatusCreated {
		t.Fatalf("create delegation: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var delegation model.Delegation
	json.NewDecoder(w.Body).Decode(&delegation)
	if delegation.ToBusinessID != biz.ID {
		t.Fatalf("expected to_business_id %s, got %s", biz.ID, delegation.ToBusinessID)
	}
	if string(delegation.Scope) != "delivery" {
		t.Fatalf("expected scope delivery, got %s", delegation.Scope)
	}

	// List delegations for share
	w = env.request("GET", "/api/v1/shares/"+shareResp.Share.ID+"/delegations", nil, userToken)
	if w.Code != http.StatusOK {
		t.Fatalf("list delegations: expected 200, got %d", w.Code)
	}
	var delegations []model.Delegation
	json.NewDecoder(w.Body).Decode(&delegations)
	if len(delegations) != 1 {
		t.Fatalf("expected 1 delegation, got %d", len(delegations))
	}

	// Revoke delegation
	w = env.request("PATCH", "/api/v1/delegations/"+delegation.ID+"/revoke", nil, userToken)
	if w.Code != http.StatusOK {
		t.Fatalf("revoke delegation: expected 200, got %d", w.Code)
	}

	// Verify revoked
	w = env.request("GET", "/api/v1/shares/"+shareResp.Share.ID+"/delegations", nil, userToken)
	json.NewDecoder(w.Body).Decode(&delegations)
	if delegations[0].Active {
		t.Fatal("expected delegation to be inactive")
	}
}

func TestDelegationScopeEnforcement(t *testing.T) {
	env := setupTestEnv(t)
	userToken, _ := env.registerUser(t, "scope-enforce@example.com", "password123", "Enforcer")

	// Create address + share with "zone" scope
	w := env.request("POST", "/api/v1/addresses", model.CreateAddressRequest{
		Label: "Office", Line1: "Zone St", City: "Vienna", PostCode: "1010", Country: "AT",
	}, userToken)
	var addr model.Address
	json.NewDecoder(w.Body).Decode(&addr)

	w = env.request("POST", "/api/v1/shares", model.CreateShareRequest{
		AddressID: addr.ID, AccessType: model.ShareAccessPublic, Scope: model.ScopeZone,
	}, userToken)
	var shareResp struct{ Share model.Share `json:"share"` }
	json.NewDecoder(w.Body).Decode(&shareResp)

	// Create business
	w = env.request("POST", "/api/v1/businesses", model.CreateBusinessRequest{Name: "Test Corp"}, userToken)
	var biz model.Business
	json.NewDecoder(w.Body).Decode(&biz)

	// Try to delegate with "full" scope — should fail (share is only "zone")
	w = env.request("POST", "/api/v1/delegations", model.CreateDelegationRequest{
		ShareID: shareResp.Share.ID, ToBusinessID: biz.ID, Scope: model.ScopeFull,
	}, userToken)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("scope upgrade: expected 400, got %d: %s", w.Code, w.Body.String())
	}

	// Delegate with "zone" scope — should work
	w = env.request("POST", "/api/v1/delegations", model.CreateDelegationRequest{
		ShareID: shareResp.Share.ID, ToBusinessID: biz.ID, Scope: model.ScopeZone,
	}, userToken)
	if w.Code != http.StatusCreated {
		t.Fatalf("zone delegation: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	// Delegate with "verify" scope — should work (lower than zone)
	w = env.request("POST", "/api/v1/delegations", model.CreateDelegationRequest{
		ShareID: shareResp.Share.ID, ToBusinessID: biz.ID, Scope: model.ScopeVerify,
	}, userToken)
	if w.Code != http.StatusCreated {
		t.Fatalf("verify delegation: expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

// ─── Phase 3: Label Tests ───────────────────────────────────

func TestLabelCreation(t *testing.T) {
	env := setupTestEnv(t)
	token, _ := env.registerUser(t, "label@example.com", "password123", "Label Test")

	// Create address + share
	w := env.request("POST", "/api/v1/addresses", model.CreateAddressRequest{
		Label: "Home", Line1: "Label Lane 42", City: "Istanbul", PostCode: "34000", Country: "TR",
	}, token)
	var addr model.Address
	json.NewDecoder(w.Body).Decode(&addr)

	w = env.request("POST", "/api/v1/shares", model.CreateShareRequest{
		AddressID: addr.ID, AccessType: model.ShareAccessPublic,
	}, token)
	var shareResp struct{ Share model.Share `json:"share"` }
	json.NewDecoder(w.Body).Decode(&shareResp)

	// Create label
	w = env.request("POST", "/api/v1/labels", model.CreateLabelRequest{
		ShareID: shareResp.Share.ID,
	}, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create label: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var labelResp model.LabelResponse
	json.NewDecoder(w.Body).Decode(&labelResp)

	if labelResp.Label.ReferenceCode == "" {
		t.Fatal("expected reference code")
	}
	if labelResp.Label.ZoneCode == "" {
		t.Fatal("expected zone code")
	}
	// Zone code should be "TR-IST-340"
	if labelResp.Label.ZoneCode != "TR-IST-340" {
		t.Fatalf("expected zone code TR-IST-340, got %s", labelResp.Label.ZoneCode)
	}
	if labelResp.QRCodeURL == "" {
		t.Fatal("expected QR code URL")
	}

	// Get label image
	w = env.request("GET", "/api/v1/labels/"+labelResp.Label.ReferenceCode+"/image", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("label image: expected 200, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "image/png" {
		t.Fatalf("expected image/png, got %s", w.Header().Get("Content-Type"))
	}
	// Check reference and zone in headers
	if w.Header().Get("X-Label-Reference") != labelResp.Label.ReferenceCode {
		t.Fatalf("expected reference in header")
	}
	if w.Header().Get("X-Label-Zone") != "TR-IST-340" {
		t.Fatalf("expected zone in header, got %s", w.Header().Get("X-Label-Zone"))
	}
}

// ─── Phase 2.5: OAuth Authorization Code Flow Tests ─────────

func TestOAuthAuthorizationCodeFlow(t *testing.T) {
	env := setupTestEnv(t)
	userToken, _ := env.registerUser(t, "oauth-user@example.com", "password123", "OAuth User")

	// Create address
	w := env.request("POST", "/api/v1/addresses", model.CreateAddressRequest{
		Label: "Home", Line1: "1 OAuth St", City: "Amsterdam", PostCode: "1012", Country: "NL", Phone: "+31123",
	}, userToken)
	var addr model.Address
	json.NewDecoder(w.Body).Decode(&addr)

	// Create business
	w = env.request("POST", "/api/v1/businesses", model.CreateBusinessRequest{Name: "ShopMart"}, userToken)
	var biz model.Business
	json.NewDecoder(w.Body).Decode(&biz)

	// Create OAuth app
	w = env.request("POST", "/api/v1/businesses/"+biz.ID+"/oauth-apps", model.CreateOAuthAppRequest{
		Name:         "ShopMart Checkout",
		RedirectURIs: []string{"https://shopmart.com/callback"},
	}, userToken)
	if w.Code != http.StatusCreated {
		t.Fatalf("create oauth app: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var appResp model.CreateOAuthAppResponse
	json.NewDecoder(w.Body).Decode(&appResp)
	if appResp.ClientSecret == "" {
		t.Fatal("expected client_secret")
	}

	// Step 1: Get authorize data (consent screen info)
	w = env.request("GET", "/api/v1/oauth/authorize?client_id="+appResp.App.ClientID+"&redirect_uri=https://shopmart.com/callback&scope=delivery&state=xyz123", nil, userToken)
	if w.Code != http.StatusOK {
		t.Fatalf("authorize: expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var consentData model.ConsentPageData
	json.NewDecoder(w.Body).Decode(&consentData)
	if consentData.App.Name != "ShopMart Checkout" {
		t.Fatalf("expected app name, got %s", consentData.App.Name)
	}
	if len(consentData.Addresses) != 1 {
		t.Fatalf("expected 1 address, got %d", len(consentData.Addresses))
	}

	// Step 2: User consents
	w = env.request("POST", "/api/v1/oauth/consent", model.ConsentRequest{
		ClientID:    appResp.App.ClientID,
		RedirectURI: "https://shopmart.com/callback",
		Scope:       "delivery",
		State:       "xyz123",
		AddressID:   addr.ID,
	}, userToken)
	if w.Code != http.StatusOK {
		t.Fatalf("consent: expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var consentResp struct {
		RedirectURL string `json:"redirect_url"`
		Code        string `json:"code"`
	}
	json.NewDecoder(w.Body).Decode(&consentResp)
	if consentResp.Code == "" {
		t.Fatal("expected authorization code")
	}

	// Step 3: Exchange code for token
	w = env.request("POST", "/api/v1/oauth/exchange", model.TokenExchangeRequest{
		GrantType:    "authorization_code",
		Code:         consentResp.Code,
		ClientID:     appResp.App.ClientID,
		ClientSecret: appResp.ClientSecret,
		RedirectURI:  "https://shopmart.com/callback",
	}, "")
	if w.Code != http.StatusOK {
		t.Fatalf("exchange: expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var tokenResp model.TokenExchangeResponse
	json.NewDecoder(w.Body).Decode(&tokenResp)
	if tokenResp.AccessToken == "" {
		t.Fatal("expected access_token")
	}
	if tokenResp.ShareToken == "" {
		t.Fatal("expected share_token")
	}
	if tokenResp.Scope != "delivery" {
		t.Fatalf("expected scope delivery, got %s", tokenResp.Scope)
	}

	// Step 4: Resolve the share token (business uses it to get the address)
	w = env.request("GET", "/api/v1/resolve/"+tokenResp.ShareToken, nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("resolve: expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resolveResp model.ResolveResponse
	json.NewDecoder(w.Body).Decode(&resolveResp)
	if resolveResp.Address.Line1 != "1 OAuth St" {
		t.Fatalf("expected '1 OAuth St', got '%s'", resolveResp.Address.Line1)
	}
	// Delivery scope: no phone
	if resolveResp.Address.Phone != "" {
		t.Fatalf("delivery scope should not include phone, got '%s'", resolveResp.Address.Phone)
	}
}

func TestOAuthInvalidRedirectURI(t *testing.T) {
	env := setupTestEnv(t)
	userToken, _ := env.registerUser(t, "oauth-bad@example.com", "password123", "Bad OAuth")

	w := env.request("POST", "/api/v1/businesses", model.CreateBusinessRequest{Name: "BadApp"}, userToken)
	var biz model.Business
	json.NewDecoder(w.Body).Decode(&biz)

	w = env.request("POST", "/api/v1/businesses/"+biz.ID+"/oauth-apps", model.CreateOAuthAppRequest{
		Name: "BadApp", RedirectURIs: []string{"https://good.com/callback"},
	}, userToken)
	var appResp model.CreateOAuthAppResponse
	json.NewDecoder(w.Body).Decode(&appResp)

	// Try authorize with wrong redirect URI
	w = env.request("GET", "/api/v1/oauth/authorize?client_id="+appResp.App.ClientID+"&redirect_uri=https://evil.com/steal&scope=full", nil, userToken)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("bad redirect: expected 400, got %d", w.Code)
	}
}

func TestOAuthCodeReuse(t *testing.T) {
	env := setupTestEnv(t)
	userToken, _ := env.registerUser(t, "oauth-reuse@example.com", "password123", "Reuse")

	w := env.request("POST", "/api/v1/addresses", model.CreateAddressRequest{
		Label: "Home", Line1: "Reuse St", City: "Rome", PostCode: "00100", Country: "IT",
	}, userToken)
	var addr model.Address
	json.NewDecoder(w.Body).Decode(&addr)

	w = env.request("POST", "/api/v1/businesses", model.CreateBusinessRequest{Name: "ReuseApp"}, userToken)
	var biz model.Business
	json.NewDecoder(w.Body).Decode(&biz)

	w = env.request("POST", "/api/v1/businesses/"+biz.ID+"/oauth-apps", model.CreateOAuthAppRequest{
		Name: "ReuseApp", RedirectURIs: []string{"https://reuse.com/cb"},
	}, userToken)
	var appResp model.CreateOAuthAppResponse
	json.NewDecoder(w.Body).Decode(&appResp)

	// Get code
	w = env.request("POST", "/api/v1/oauth/consent", model.ConsentRequest{
		ClientID: appResp.App.ClientID, RedirectURI: "https://reuse.com/cb", Scope: "full", AddressID: addr.ID,
	}, userToken)
	var consentResp struct{ Code string `json:"code"` }
	json.NewDecoder(w.Body).Decode(&consentResp)

	// First exchange — OK
	w = env.request("POST", "/api/v1/oauth/exchange", model.TokenExchangeRequest{
		GrantType: "authorization_code", Code: consentResp.Code,
		ClientID: appResp.App.ClientID, ClientSecret: appResp.ClientSecret, RedirectURI: "https://reuse.com/cb",
	}, "")
	if w.Code != http.StatusOK {
		t.Fatalf("first exchange: expected 200, got %d", w.Code)
	}

	// Second exchange — should fail (code already used)
	w = env.request("POST", "/api/v1/oauth/exchange", model.TokenExchangeRequest{
		GrantType: "authorization_code", Code: consentResp.Code,
		ClientID: appResp.App.ClientID, ClientSecret: appResp.ClientSecret, RedirectURI: "https://reuse.com/cb",
	}, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("code reuse: expected 400, got %d", w.Code)
	}
}
