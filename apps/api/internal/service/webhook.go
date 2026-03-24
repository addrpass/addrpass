package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/addrpass/addrpass/apps/api/internal/model"
)

var ErrWebhookNotFound = errors.New("webhook not found")

type WebhookService struct {
	db     *pgxpool.Pool
	client *http.Client
}

func NewWebhookService(db *pgxpool.Pool) *WebhookService {
	return &WebhookService{
		db: db,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *WebhookService) Create(ctx context.Context, userID string, req model.CreateWebhookRequest) (*model.Webhook, error) {
	events := req.Events
	if len(events) == 0 {
		events = []string{"access"}
	}

	// Generate webhook signing secret
	secret, err := generateToken()
	if err != nil {
		return nil, err
	}

	var wh model.Webhook
	err = s.db.QueryRow(ctx,
		`INSERT INTO webhooks (user_id, url, secret, events)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, url, secret, events, active, failure_count, last_triggered_at, created_at`,
		userID, req.URL, secret, events,
	).Scan(&wh.ID, &wh.UserID, &wh.URL, &wh.Secret, &wh.Events, &wh.Active, &wh.FailureCount, &wh.LastTriggeredAt, &wh.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &wh, nil
}

func (s *WebhookService) List(ctx context.Context, userID string) ([]model.Webhook, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, user_id, url, secret, events, active, failure_count, last_triggered_at, created_at
		 FROM webhooks WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []model.Webhook
	for rows.Next() {
		var wh model.Webhook
		if err := rows.Scan(&wh.ID, &wh.UserID, &wh.URL, &wh.Secret, &wh.Events, &wh.Active, &wh.FailureCount, &wh.LastTriggeredAt, &wh.CreatedAt); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, wh)
	}
	if webhooks == nil {
		webhooks = []model.Webhook{}
	}
	return webhooks, nil
}

func (s *WebhookService) Delete(ctx context.Context, userID, webhookID string) error {
	result, err := s.db.Exec(ctx, `DELETE FROM webhooks WHERE id = $1 AND user_id = $2`, webhookID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrWebhookNotFound
	}
	return nil
}

// DispatchAccessEvent fires webhooks for the share owner when their address is accessed.
// Called asynchronously (in a goroutine) from the resolve handler.
func (s *WebhookService) DispatchAccessEvent(shareOwnerID string, event model.AccessEvent) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := s.db.Query(ctx,
		`SELECT id, url, secret FROM webhooks
		 WHERE user_id = $1 AND active = TRUE AND 'access' = ANY(events)`,
		shareOwnerID,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	payload := model.WebhookEvent{
		Event:     "access",
		Timestamp: time.Now().UTC(),
		Data:      event,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return
	}

	for rows.Next() {
		var whID, url, secret string
		if err := rows.Scan(&whID, &url, &secret); err != nil {
			continue
		}
		go s.deliver(whID, url, secret, "access", body)
	}
}

func (s *WebhookService) deliver(webhookID, url, secret, event string, body []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		s.recordDelivery(webhookID, event, body, 0, err.Error(), false)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-AddrPass-Event", event)
	if secret != "" {
		sig := computeSignature(body, secret)
		req.Header.Set("X-AddrPass-Signature", sig)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.recordDelivery(webhookID, event, body, 0, err.Error(), false)
		s.incrementFailure(webhookID)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	success := resp.StatusCode >= 200 && resp.StatusCode < 300

	s.recordDelivery(webhookID, event, body, resp.StatusCode, string(respBody), success)

	if success {
		s.resetFailure(webhookID)
	} else {
		s.incrementFailure(webhookID)
	}
}

func (s *WebhookService) recordDelivery(webhookID, event string, payload []byte, statusCode int, responseBody string, success bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = s.db.Exec(ctx,
		`INSERT INTO webhook_deliveries (webhook_id, event, payload, status_code, response_body, success)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		webhookID, event, payload, statusCode, responseBody, success,
	)
	_, _ = s.db.Exec(ctx, `UPDATE webhooks SET last_triggered_at = NOW() WHERE id = $1`, webhookID)
}

func (s *WebhookService) incrementFailure(webhookID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Auto-disable after 10 consecutive failures
	_, _ = s.db.Exec(ctx,
		`UPDATE webhooks SET failure_count = failure_count + 1,
		 active = CASE WHEN failure_count + 1 >= 10 THEN FALSE ELSE active END
		 WHERE id = $1`, webhookID)
}

func (s *WebhookService) resetFailure(webhookID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = s.db.Exec(ctx, `UPDATE webhooks SET failure_count = 0 WHERE id = $1`, webhookID)
}

func computeSignature(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// GetByShareOwner finds the share owner's user_id for a given share_id
func (s *WebhookService) GetShareOwnerID(ctx context.Context, shareID string) (string, error) {
	var userID string
	err := s.db.QueryRow(ctx, `SELECT user_id FROM shares WHERE id = $1`, shareID).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrShareNotFound
		}
		return "", err
	}
	return userID, nil
}
