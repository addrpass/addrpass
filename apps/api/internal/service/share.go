package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/addrpass/addrpass/apps/api/internal/model"
)

var (
	ErrShareNotFound = errors.New("share not found")
	ErrShareExpired  = errors.New("share has expired")
	ErrShareRevoked  = errors.New("share has been revoked")
	ErrMaxAccesses   = errors.New("maximum accesses reached")
	ErrInvalidPin    = errors.New("invalid pin")
)

type ShareService struct {
	db *pgxpool.Pool
}

func NewShareService(db *pgxpool.Pool) *ShareService {
	return &ShareService{db: db}
}

func (s *ShareService) Create(ctx context.Context, userID string, req model.CreateShareRequest) (*model.Share, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	if req.AccessType == "" {
		req.AccessType = model.ShareAccessPublic
	}

	var share model.Share
	err = s.db.QueryRow(ctx,
		`INSERT INTO shares (address_id, user_id, token, access_type, pin, expires_at, max_accesses)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, address_id, user_id, token, access_type, pin, expires_at, max_accesses, access_count, active, created_at`,
		req.AddressID, userID, token, req.AccessType, req.Pin, req.ExpiresAt, req.MaxAccesses,
	).Scan(&share.ID, &share.AddressID, &share.UserID, &share.Token, &share.AccessType, &share.Pin, &share.ExpiresAt, &share.MaxAccesses, &share.AccessCount, &share.Active, &share.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &share, nil
}

func (s *ShareService) List(ctx context.Context, userID string) ([]model.Share, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, address_id, user_id, token, access_type, pin, expires_at, max_accesses, access_count, active, created_at
		 FROM shares WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shares []model.Share
	for rows.Next() {
		var sh model.Share
		if err := rows.Scan(&sh.ID, &sh.AddressID, &sh.UserID, &sh.Token, &sh.AccessType, &sh.Pin, &sh.ExpiresAt, &sh.MaxAccesses, &sh.AccessCount, &sh.Active, &sh.CreatedAt); err != nil {
			return nil, err
		}
		shares = append(shares, sh)
	}
	if shares == nil {
		shares = []model.Share{}
	}
	return shares, nil
}

func (s *ShareService) Revoke(ctx context.Context, userID, shareID string) error {
	result, err := s.db.Exec(ctx,
		`UPDATE shares SET active = FALSE WHERE id = $1 AND user_id = $2`,
		shareID, userID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrShareNotFound
	}
	return nil
}

func (s *ShareService) Delete(ctx context.Context, userID, shareID string) error {
	result, err := s.db.Exec(ctx,
		`DELETE FROM shares WHERE id = $1 AND user_id = $2`,
		shareID, userID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrShareNotFound
	}
	return nil
}

func (s *ShareService) Resolve(ctx context.Context, token, pin, ip, userAgent string) (*model.Address, error) {
	var share model.Share
	err := s.db.QueryRow(ctx,
		`SELECT id, address_id, user_id, token, access_type, pin, expires_at, max_accesses, access_count, active
		 FROM shares WHERE token = $1`,
		token,
	).Scan(&share.ID, &share.AddressID, &share.UserID, &share.Token, &share.AccessType, &share.Pin, &share.ExpiresAt, &share.MaxAccesses, &share.AccessCount, &share.Active)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrShareNotFound
		}
		return nil, err
	}

	if !share.Active {
		return nil, ErrShareRevoked
	}
	if share.ExpiresAt != nil && time.Now().After(*share.ExpiresAt) {
		return nil, ErrShareExpired
	}
	if share.MaxAccesses != nil && share.AccessCount >= *share.MaxAccesses {
		return nil, ErrMaxAccesses
	}
	if share.Pin != "" && pin != share.Pin {
		return nil, ErrInvalidPin
	}

	// Fetch the address
	var addr model.Address
	err = s.db.QueryRow(ctx,
		`SELECT id, user_id, label, line1, line2, city, state, post_code, country, phone, created_at, updated_at
		 FROM addresses WHERE id = $1`,
		share.AddressID,
	).Scan(&addr.ID, &addr.UserID, &addr.Label, &addr.Line1, &addr.Line2, &addr.City, &addr.State, &addr.PostCode, &addr.Country, &addr.Phone, &addr.CreatedAt, &addr.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Increment access count and log
	_, _ = s.db.Exec(ctx, `UPDATE shares SET access_count = access_count + 1 WHERE id = $1`, share.ID)
	_, _ = s.db.Exec(ctx,
		`INSERT INTO access_logs (share_id, ip, user_agent) VALUES ($1, $2, $3)`,
		share.ID, ip, userAgent,
	)

	return &addr, nil
}

func (s *ShareService) GetAccessLogs(ctx context.Context, userID, shareID string) ([]model.AccessLog, error) {
	// Verify ownership
	var count int
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM shares WHERE id = $1 AND user_id = $2`, shareID, userID).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ErrShareNotFound
	}

	rows, err := s.db.Query(ctx,
		`SELECT id, share_id, ip, user_agent, country, access_at
		 FROM access_logs WHERE share_id = $1 ORDER BY access_at DESC LIMIT 100`,
		shareID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.AccessLog
	for rows.Next() {
		var l model.AccessLog
		if err := rows.Scan(&l.ID, &l.ShareID, &l.IP, &l.UserAgent, &l.Country, &l.AccessAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	if logs == nil {
		logs = []model.AccessLog{}
	}
	return logs, nil
}

func generateToken() (string, error) {
	b := make([]byte, 18)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}
