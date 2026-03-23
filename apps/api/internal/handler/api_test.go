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
	pool.Exec(ctx, "DELETE FROM access_logs")
	pool.Exec(ctx, "DELETE FROM shares")
	pool.Exec(ctx, "DELETE FROM addresses")
	pool.Exec(ctx, "DELETE FROM users")

	authSvc := service.NewAuthService(pool, testJWTSecret)
	addressSvc := service.NewAddressService(pool)
	shareSvc := service.NewShareService(pool)

	authH := handler.NewAuthHandler(authSvc)
	addressH := handler.NewAddressHandler(addressSvc)
	shareH := handler.NewShareHandler(shareSvc, testBaseURL)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"*"}}))

	r.Post("/api/v1/auth/register", authH.Register)
	r.Post("/api/v1/auth/login", authH.Login)
	r.Get("/api/v1/resolve/{token}", shareH.Resolve)
	r.Get("/api/v1/qr/{token}", shareH.QRCode)

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
	})

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
