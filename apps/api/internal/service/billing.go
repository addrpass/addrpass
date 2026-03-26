package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Plan limits
var PlanLimits = map[string]PlanConfig{
	"free":       {Addresses: 3, Shares: 10, ResolutionsPerMonth: 50, HardLimit: true},
	"pro":        {Addresses: -1, Shares: -1, ResolutionsPerMonth: 1000, HardLimit: false},
	"business":   {Addresses: -1, Shares: -1, ResolutionsPerMonth: 10000, HardLimit: false},
	"enterprise": {Addresses: -1, Shares: -1, ResolutionsPerMonth: -1, HardLimit: false},
}

type PlanConfig struct {
	Addresses           int  // -1 = unlimited
	Shares              int  // -1 = unlimited
	ResolutionsPerMonth int  // -1 = unlimited
	HardLimit           bool // true = reject when over limit; false = allow and bill overage
}

type BillingService struct {
	db *pgxpool.Pool
}

func NewBillingService(db *pgxpool.Pool) *BillingService {
	return &BillingService{db: db}
}

// GetUserPlan returns the user's current plan.
func (s *BillingService) GetUserPlan(ctx context.Context, userID string) (string, error) {
	var plan string
	err := s.db.QueryRow(ctx, `SELECT plan FROM users WHERE id = $1`, userID).Scan(&plan)
	if err != nil {
		return "free", err
	}
	return plan, nil
}

// CurrentMonth returns the current billing month string (e.g., "2026-03").
func CurrentMonth() string {
	return time.Now().UTC().Format("2006-01")
}

// IncrementResolutions atomically increments the resolution counter for the current month.
// Returns the new count and whether the resolution should be allowed.
func (s *BillingService) IncrementResolutions(ctx context.Context, userID string) (count int, allowed bool, err error) {
	month := CurrentMonth()

	// Upsert usage record and increment
	err = s.db.QueryRow(ctx,
		`INSERT INTO usage_records (user_id, month, resolutions)
		 VALUES ($1, $2, 1)
		 ON CONFLICT (user_id, month)
		 DO UPDATE SET resolutions = usage_records.resolutions + 1
		 RETURNING resolutions`,
		userID, month,
	).Scan(&count)
	if err != nil {
		return 0, false, err
	}

	// Check plan limits
	var plan string
	err = s.db.QueryRow(ctx, `SELECT plan FROM users WHERE id = $1`, userID).Scan(&plan)
	if err != nil {
		return count, true, nil // Allow on error (don't block service)
	}

	config, ok := PlanLimits[plan]
	if !ok {
		return count, true, nil
	}

	// Unlimited plan
	if config.ResolutionsPerMonth < 0 {
		return count, true, nil
	}

	// Within limit
	if count <= config.ResolutionsPerMonth {
		return count, true, nil
	}

	// Over limit — hard limit rejects, soft limit allows (bills overage)
	if config.HardLimit {
		return count, false, nil
	}

	return count, true, nil
}

// GetMonthlyUsage returns the current month's resolution count for a user.
func (s *BillingService) GetMonthlyUsage(ctx context.Context, userID string) (int, error) {
	month := CurrentMonth()
	var count int
	err := s.db.QueryRow(ctx,
		`SELECT COALESCE((SELECT resolutions FROM usage_records WHERE user_id = $1 AND month = $2), 0)`,
		userID, month,
	).Scan(&count)
	return count, err
}

// GetUsageHistory returns the last N months of usage for a user.
func (s *BillingService) GetUsageHistory(ctx context.Context, userID string, months int) ([]UsageRecord, error) {
	rows, err := s.db.Query(ctx,
		`SELECT month, resolutions FROM usage_records
		 WHERE user_id = $1 ORDER BY month DESC LIMIT $2`,
		userID, months,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []UsageRecord
	for rows.Next() {
		var r UsageRecord
		if err := rows.Scan(&r.Month, &r.Resolutions); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	if records == nil {
		records = []UsageRecord{}
	}
	return records, nil
}

type UsageRecord struct {
	Month       string `json:"month"`
	Resolutions int    `json:"resolutions"`
}

// CheckAddressLimit checks if the user can create another address.
func (s *BillingService) CheckAddressLimit(ctx context.Context, userID string) (bool, error) {
	var plan string
	var count int
	err := s.db.QueryRow(ctx, `SELECT plan FROM users WHERE id = $1`, userID).Scan(&plan)
	if err != nil {
		return true, nil
	}

	config, ok := PlanLimits[plan]
	if !ok || config.Addresses < 0 {
		return true, nil
	}

	err = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM addresses WHERE user_id = $1`, userID).Scan(&count)
	if err != nil {
		return true, nil
	}

	return count < config.Addresses, nil
}

// CheckShareLimit checks if the user can create another share.
func (s *BillingService) CheckShareLimit(ctx context.Context, userID string) (bool, error) {
	var plan string
	var count int
	err := s.db.QueryRow(ctx, `SELECT plan FROM users WHERE id = $1`, userID).Scan(&plan)
	if err != nil {
		return true, nil
	}

	config, ok := PlanLimits[plan]
	if !ok || config.Shares < 0 {
		return true, nil
	}

	err = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM shares WHERE user_id = $1 AND active = TRUE`, userID).Scan(&count)
	if err != nil {
		return true, nil
	}

	return count < config.Shares, nil
}

// UpdatePlan updates a user's plan and resolution limit.
func (s *BillingService) UpdatePlan(ctx context.Context, userID, plan string) error {
	config, ok := PlanLimits[plan]
	if !ok {
		return fmt.Errorf("unknown plan: %s", plan)
	}

	limit := config.ResolutionsPerMonth
	if limit < 0 {
		limit = 999999999
	}

	_, err := s.db.Exec(ctx,
		`UPDATE users SET plan = $1, plan_resolution_limit = $2 WHERE id = $3`,
		plan, limit, userID,
	)
	return err
}

// SetStripeCustomer links a Stripe customer ID to a user.
func (s *BillingService) SetStripeCustomer(ctx context.Context, userID, customerID string) error {
	_, err := s.db.Exec(ctx,
		`UPDATE users SET stripe_customer_id = $1 WHERE id = $2`,
		customerID, userID,
	)
	return err
}

// SetStripeSubscription links a Stripe subscription ID to a user.
func (s *BillingService) SetStripeSubscription(ctx context.Context, userID, subscriptionID string) error {
	_, err := s.db.Exec(ctx,
		`UPDATE users SET stripe_subscription_id = $1 WHERE id = $2`,
		subscriptionID, userID,
	)
	return err
}

// LogBillingEvent records a billing event for audit.
func (s *BillingService) LogBillingEvent(ctx context.Context, userID, eventID, eventType string, payload interface{}) {
	_, _ = s.db.Exec(ctx,
		`INSERT INTO billing_events (user_id, stripe_event_id, event_type, payload)
		 VALUES ($1, $2, $3, $4)`,
		nilIfEmpty(userID), eventID, eventType, payload,
	)
}

// GetUnreportedUsage returns all usage records that haven't been reported to Stripe yet.
func (s *BillingService) GetUnreportedUsage(ctx context.Context) ([]UnreportedUsage, error) {
	rows, err := s.db.Query(ctx,
		`SELECT ur.id, ur.user_id, ur.month, ur.resolutions, u.plan, u.stripe_subscription_id
		 FROM usage_records ur
		 JOIN users u ON ur.user_id = u.id
		 WHERE ur.reported_to_stripe = FALSE
		   AND u.plan IN ('pro', 'business')
		   AND u.stripe_subscription_id IS NOT NULL
		   AND ur.month < $1`,
		CurrentMonth(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []UnreportedUsage
	for rows.Next() {
		var r UnreportedUsage
		if err := rows.Scan(&r.RecordID, &r.UserID, &r.Month, &r.Resolutions, &r.Plan, &r.SubscriptionID); err != nil {
			return nil, err
		}

		config := PlanLimits[r.Plan]
		overage := r.Resolutions - config.ResolutionsPerMonth
		if overage > 0 {
			r.Overage = overage
		}
		records = append(records, r)
	}
	if records == nil {
		records = []UnreportedUsage{}
	}
	return records, nil
}

type UnreportedUsage struct {
	RecordID       string `json:"record_id"`
	UserID         string `json:"user_id"`
	Month          string `json:"month"`
	Resolutions    int    `json:"resolutions"`
	Plan           string `json:"plan"`
	SubscriptionID string `json:"subscription_id"`
	Overage        int    `json:"overage"`
}

// MarkReported marks a usage record as reported to Stripe.
func (s *BillingService) MarkReported(ctx context.Context, recordID string) error {
	_, err := s.db.Exec(ctx, `UPDATE usage_records SET reported_to_stripe = TRUE WHERE id = $1`, recordID)
	return err
}
