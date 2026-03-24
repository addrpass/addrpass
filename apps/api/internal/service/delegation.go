package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/addrpass/addrpass/apps/api/internal/model"
)

var (
	ErrDelegationNotFound = errors.New("delegation not found")
	ErrDelegationExpired  = errors.New("delegation has expired")
	ErrDelegationRevoked  = errors.New("delegation has been revoked")
	ErrDelegationMaxHits  = errors.New("delegation max accesses reached")
	ErrNotDelegated       = errors.New("business is not delegated for this share")
	ErrCannotDelegateUp   = errors.New("cannot delegate higher scope than you have")
)

// Scope hierarchy: full > delivery > zone > verify
var scopeLevel = map[model.ShareScope]int{
	model.ScopeFull:     4,
	model.ScopeDelivery: 3,
	model.ScopeZone:     2,
	model.ScopeVerify:   1,
}

type DelegationService struct {
	db *pgxpool.Pool
}

func NewDelegationService(db *pgxpool.Pool) *DelegationService {
	return &DelegationService{db: db}
}

// CreateByUser creates a delegation from the share owner to a business.
// The user must own the share.
func (s *DelegationService) CreateByUser(ctx context.Context, userID string, req model.CreateDelegationRequest) (*model.Delegation, error) {
	// Verify user owns the share
	var shareScope model.ShareScope
	var count int
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*), scope FROM shares WHERE id = $1 AND user_id = $2 AND active = TRUE GROUP BY scope`,
		req.ShareID, userID,
	).Scan(&count, &shareScope)
	if err != nil || count == 0 {
		return nil, ErrShareNotFound
	}

	scope := req.Scope
	if scope == "" {
		scope = shareScope
	}

	// Cannot delegate higher scope than the share itself
	if scopeLevel[scope] > scopeLevel[shareScope] {
		return nil, ErrCannotDelegateUp
	}

	return s.insertDelegation(ctx, req.ShareID, "", req.ToBusinessID, scope, req.ExpiresAt, req.MaxAccesses, req.Note)
}

// CreateByBusiness creates a sub-delegation from one business to another.
// The from_business must already have a valid delegation for this share.
func (s *DelegationService) CreateByBusiness(ctx context.Context, fromBusinessID string, req model.CreateDelegationRequest) (*model.Delegation, error) {
	// Verify the from_business has a valid delegation
	parentScope, err := s.getEffectiveScope(ctx, req.ShareID, fromBusinessID)
	if err != nil {
		return nil, err
	}

	scope := req.Scope
	if scope == "" {
		scope = parentScope
	}

	if scopeLevel[scope] > scopeLevel[parentScope] {
		return nil, ErrCannotDelegateUp
	}

	return s.insertDelegation(ctx, req.ShareID, fromBusinessID, req.ToBusinessID, scope, req.ExpiresAt, req.MaxAccesses, req.Note)
}

func (s *DelegationService) insertDelegation(ctx context.Context, shareID, fromBizID, toBizID string, scope model.ShareScope, expiresAt *time.Time, maxAccesses *int, note string) (*model.Delegation, error) {
	var d model.Delegation
	err := s.db.QueryRow(ctx,
		`INSERT INTO delegations (share_id, from_business_id, to_business_id, scope, expires_at, max_accesses, note)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, share_id, COALESCE(from_business_id::text, ''), to_business_id, scope, expires_at, max_accesses, access_count, active, note, created_at`,
		shareID, nilIfEmpty(fromBizID), toBizID, scope, expiresAt, maxAccesses, note,
	).Scan(&d.ID, &d.ShareID, &d.FromBusinessID, &d.ToBusinessID, &d.Scope, &d.ExpiresAt, &d.MaxAccesses, &d.AccessCount, &d.Active, &d.Note, &d.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// ListForShare lists all delegations for a share (owner view).
func (s *DelegationService) ListForShare(ctx context.Context, userID, shareID string) ([]model.Delegation, error) {
	// Verify ownership
	var count int
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM shares WHERE id = $1 AND user_id = $2`, shareID, userID).Scan(&count)
	if err != nil || count == 0 {
		return nil, ErrShareNotFound
	}

	return s.queryDelegations(ctx, `SELECT id, share_id, COALESCE(from_business_id::text, ''), to_business_id, scope, expires_at, max_accesses, access_count, active, note, created_at FROM delegations WHERE share_id = $1 ORDER BY created_at DESC`, shareID)
}

// Revoke deactivates a delegation. Only the share owner can do this.
func (s *DelegationService) Revoke(ctx context.Context, userID, delegationID string) error {
	result, err := s.db.Exec(ctx,
		`UPDATE delegations SET active = FALSE
		 WHERE id = $1 AND share_id IN (SELECT id FROM shares WHERE user_id = $2)`,
		delegationID, userID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrDelegationNotFound
	}
	return nil
}

// getEffectiveScope returns the scope a business has for a given share,
// either through direct delegation or the share's delegated_to_business_id.
func (s *DelegationService) getEffectiveScope(ctx context.Context, shareID, businessID string) (model.ShareScope, error) {
	var scope model.ShareScope
	var active bool
	var expiresAt *time.Time
	var maxAccesses *int
	var accessCount int

	err := s.db.QueryRow(ctx,
		`SELECT scope, active, expires_at, max_accesses, access_count FROM delegations
		 WHERE share_id = $1 AND to_business_id = $2 AND active = TRUE
		 ORDER BY created_at DESC LIMIT 1`,
		shareID, businessID,
	).Scan(&scope, &active, &expiresAt, &maxAccesses, &accessCount)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotDelegated
		}
		return "", err
	}

	if !active {
		return "", ErrDelegationRevoked
	}
	if expiresAt != nil && time.Now().After(*expiresAt) {
		return "", ErrDelegationExpired
	}
	if maxAccesses != nil && accessCount >= *maxAccesses {
		return "", ErrDelegationMaxHits
	}

	return scope, nil
}

// ResolveWithDelegation resolves a share token using a business delegation.
// Returns the address scoped to the delegation's scope (not the share's scope).
func (s *DelegationService) ResolveWithDelegation(ctx context.Context, shareToken, businessID string) (model.ShareScope, string, error) {
	// Find the share
	var shareID string
	err := s.db.QueryRow(ctx, `SELECT id FROM shares WHERE token = $1 AND active = TRUE`, shareToken).Scan(&shareID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", ErrShareNotFound
		}
		return "", "", err
	}

	// Find delegation for this business
	var delegationID string
	var scope model.ShareScope
	var active bool
	var expiresAt *time.Time
	var maxAccesses *int
	var accessCount int

	err = s.db.QueryRow(ctx,
		`SELECT id, scope, active, expires_at, max_accesses, access_count FROM delegations
		 WHERE share_id = $1 AND to_business_id = $2 AND active = TRUE
		 ORDER BY created_at DESC LIMIT 1`,
		shareID, businessID,
	).Scan(&delegationID, &scope, &active, &expiresAt, &maxAccesses, &accessCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", ErrNotDelegated
		}
		return "", "", err
	}

	if expiresAt != nil && time.Now().After(*expiresAt) {
		return "", "", ErrDelegationExpired
	}
	if maxAccesses != nil && accessCount >= *maxAccesses {
		return "", "", ErrDelegationMaxHits
	}

	// Increment delegation access count
	_, _ = s.db.Exec(ctx, `UPDATE delegations SET access_count = access_count + 1 WHERE id = $1`, delegationID)

	return scope, delegationID, nil
}

func (s *DelegationService) queryDelegations(ctx context.Context, query string, args ...interface{}) ([]model.Delegation, error) {
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var delegations []model.Delegation
	for rows.Next() {
		var d model.Delegation
		if err := rows.Scan(&d.ID, &d.ShareID, &d.FromBusinessID, &d.ToBusinessID, &d.Scope, &d.ExpiresAt, &d.MaxAccesses, &d.AccessCount, &d.Active, &d.Note, &d.CreatedAt); err != nil {
			return nil, err
		}
		delegations = append(delegations, d)
	}
	if delegations == nil {
		delegations = []model.Delegation{}
	}
	return delegations, nil
}
