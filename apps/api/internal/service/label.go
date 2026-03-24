package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/addrpass/addrpass/apps/api/internal/model"
)

var ErrLabelNotFound = errors.New("label not found")

type LabelService struct {
	db *pgxpool.Pool
}

func NewLabelService(db *pgxpool.Pool) *LabelService {
	return &LabelService{db: db}
}

// Create generates a privacy-preserving shipping label for a share.
// The label contains a reference code and zone code (derived from address),
// but no plaintext address.
func (s *LabelService) Create(ctx context.Context, shareID, businessID string) (*model.Label, error) {
	// Fetch share + address for zone code generation
	var addressCity, addressPostCode, addressCountry string
	var shareActive bool
	err := s.db.QueryRow(ctx,
		`SELECT a.city, a.post_code, a.country, s.active
		 FROM shares s JOIN addresses a ON s.address_id = a.id
		 WHERE s.id = $1`,
		shareID,
	).Scan(&addressCity, &addressPostCode, &addressCountry, &shareActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrShareNotFound
		}
		return nil, err
	}
	if !shareActive {
		return nil, ErrShareRevoked
	}

	refCode, err := generateReferenceCode()
	if err != nil {
		return nil, err
	}

	zoneCode := generateZoneCode(addressCountry, addressCity, addressPostCode)

	var label model.Label
	err = s.db.QueryRow(ctx,
		`INSERT INTO labels (share_id, business_id, reference_code, zone_code)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, share_id, COALESCE(business_id::text, ''), reference_code, zone_code, format, created_at`,
		shareID, nilIfEmpty(businessID), refCode, zoneCode,
	).Scan(&label.ID, &label.ShareID, &label.BusinessID, &label.ReferenceCode, &label.ZoneCode, &label.Format, &label.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &label, nil
}

// GetByReference looks up a label by its reference code.
func (s *LabelService) GetByReference(ctx context.Context, refCode string) (*model.Label, error) {
	var label model.Label
	err := s.db.QueryRow(ctx,
		`SELECT id, share_id, COALESCE(business_id::text, ''), reference_code, zone_code, format, created_at
		 FROM labels WHERE reference_code = $1`,
		refCode,
	).Scan(&label.ID, &label.ShareID, &label.BusinessID, &label.ReferenceCode, &label.ZoneCode, &label.Format, &label.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrLabelNotFound
		}
		return nil, err
	}
	return &label, nil
}

// GetShareTokenForLabel returns the share token for a label, used for QR code generation.
func (s *LabelService) GetShareTokenForLabel(ctx context.Context, labelID string) (string, error) {
	var token string
	err := s.db.QueryRow(ctx,
		`SELECT s.token FROM labels l JOIN shares s ON l.share_id = s.id WHERE l.id = $1`,
		labelID,
	).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}

// generateReferenceCode creates a human-readable reference like "AP-7X3K-9M2P"
func generateReferenceCode() (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // no I, O, 0, 1 to avoid confusion
	parts := make([]string, 2)
	for i := range parts {
		seg := make([]byte, 4)
		for j := range seg {
			n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err != nil {
				return "", err
			}
			seg[j] = charset[n.Int64()]
		}
		parts[i] = string(seg)
	}
	return fmt.Sprintf("AP-%s-%s", parts[0], parts[1]), nil
}

// generateZoneCode creates a routing zone code from country + city + postal code.
// Format: "DE-MUC-803" (country-city_abbrev-postal_prefix)
func generateZoneCode(country, city, postCode string) string {
	country = strings.ToUpper(strings.TrimSpace(country))

	cityAbbr := strings.ToUpper(strings.TrimSpace(city))
	if len(cityAbbr) > 3 {
		cityAbbr = cityAbbr[:3]
	}

	postalPrefix := strings.TrimSpace(postCode)
	if len(postalPrefix) > 3 {
		postalPrefix = postalPrefix[:3]
	}

	return fmt.Sprintf("%s-%s-%s", country, cityAbbr, postalPrefix)
}
