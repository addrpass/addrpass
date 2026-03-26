package handler

import (
	"net/http"

	"github.com/addrpass/addrpass/apps/api/internal/middleware"
	"github.com/addrpass/addrpass/apps/api/internal/service"
)

type BillingHandler struct {
	billing *service.BillingService
}

func NewBillingHandler(billing *service.BillingService) *BillingHandler {
	return &BillingHandler{billing: billing}
}

// GetUsage returns the current user's usage and plan info.
func (h *BillingHandler) GetUsage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	plan, _ := h.billing.GetUserPlan(r.Context(), userID)
	usage, _ := h.billing.GetMonthlyUsage(r.Context(), userID)
	history, _ := h.billing.GetUsageHistory(r.Context(), userID, 6)

	config := service.PlanLimits[plan]
	limit := config.ResolutionsPerMonth
	if limit < 0 {
		limit = -1 // unlimited
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"plan":              plan,
		"resolutions_used":  usage,
		"resolutions_limit": limit,
		"hard_limit":        config.HardLimit,
		"history":           history,
		"month":             service.CurrentMonth(),
	})
}

// GetPlanLimits returns the limits for all plans (public).
func (h *BillingHandler) GetPlanLimits(w http.ResponseWriter, r *http.Request) {
	plans := map[string]interface{}{
		"free": map[string]interface{}{
			"price":        0,
			"addresses":    3,
			"shares":       10,
			"resolutions":  50,
			"hard_limit":   true,
			"features":     []string{"QR codes", "PIN protection", "Expiration", "Access logs"},
		},
		"pro": map[string]interface{}{
			"price":        9,
			"addresses":    -1,
			"shares":       -1,
			"resolutions":  1000,
			"hard_limit":   false,
			"overage_rate": 0.005,
			"features":     []string{"Everything in Free", "Webhooks", "API access", "Priority support"},
		},
		"business": map[string]interface{}{
			"price":        49,
			"addresses":    -1,
			"shares":       -1,
			"resolutions":  10000,
			"hard_limit":   false,
			"overage_rate": 0.005,
			"features":     []string{"Everything in Pro", "API keys", "OAuth2", "Delegations", "Shipping labels", "Team management"},
		},
		"enterprise": map[string]interface{}{
			"price":        -1,
			"addresses":    -1,
			"shares":       -1,
			"resolutions":  -1,
			"hard_limit":   false,
			"features":     []string{"Everything in Business", "Unlimited resolutions", "SLA", "Dedicated support", "Data residency"},
		},
	}

	writeJSON(w, http.StatusOK, plans)
}
