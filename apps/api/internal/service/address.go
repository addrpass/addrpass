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
	ErrAddressNotFound = errors.New("address not found")
	ErrNotOwner        = errors.New("not the owner of this address")
)

type AddressService struct {
	db *pgxpool.Pool
}

func NewAddressService(db *pgxpool.Pool) *AddressService {
	return &AddressService{db: db}
}

func (s *AddressService) Create(ctx context.Context, userID string, req model.CreateAddressRequest) (*model.Address, error) {
	var addr model.Address
	err := s.db.QueryRow(ctx,
		`INSERT INTO addresses (user_id, label, line1, line2, city, state, post_code, country, phone)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, user_id, label, line1, line2, city, state, post_code, country, phone, created_at, updated_at`,
		userID, req.Label, req.Line1, req.Line2, req.City, req.State, req.PostCode, req.Country, req.Phone,
	).Scan(&addr.ID, &addr.UserID, &addr.Label, &addr.Line1, &addr.Line2, &addr.City, &addr.State, &addr.PostCode, &addr.Country, &addr.Phone, &addr.CreatedAt, &addr.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

func (s *AddressService) List(ctx context.Context, userID string) ([]model.Address, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, user_id, label, line1, line2, city, state, post_code, country, phone, created_at, updated_at
		 FROM addresses WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []model.Address
	for rows.Next() {
		var a model.Address
		if err := rows.Scan(&a.ID, &a.UserID, &a.Label, &a.Line1, &a.Line2, &a.City, &a.State, &a.PostCode, &a.Country, &a.Phone, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	if addresses == nil {
		addresses = []model.Address{}
	}
	return addresses, nil
}

func (s *AddressService) Get(ctx context.Context, userID, addressID string) (*model.Address, error) {
	var addr model.Address
	err := s.db.QueryRow(ctx,
		`SELECT id, user_id, label, line1, line2, city, state, post_code, country, phone, created_at, updated_at
		 FROM addresses WHERE id = $1 AND user_id = $2`,
		addressID, userID,
	).Scan(&addr.ID, &addr.UserID, &addr.Label, &addr.Line1, &addr.Line2, &addr.City, &addr.State, &addr.PostCode, &addr.Country, &addr.Phone, &addr.CreatedAt, &addr.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAddressNotFound
		}
		return nil, err
	}
	return &addr, nil
}

func (s *AddressService) Update(ctx context.Context, userID, addressID string, req model.UpdateAddressRequest) (*model.Address, error) {
	// Verify ownership first
	if _, err := s.Get(ctx, userID, addressID); err != nil {
		return nil, err
	}

	var addr model.Address
	err := s.db.QueryRow(ctx,
		`UPDATE addresses SET
			label = COALESCE($3, label),
			line1 = COALESCE($4, line1),
			line2 = COALESCE($5, line2),
			city = COALESCE($6, city),
			state = COALESCE($7, state),
			post_code = COALESCE($8, post_code),
			country = COALESCE($9, country),
			phone = COALESCE($10, phone),
			updated_at = $11
		 WHERE id = $1 AND user_id = $2
		 RETURNING id, user_id, label, line1, line2, city, state, post_code, country, phone, created_at, updated_at`,
		addressID, userID, req.Label, req.Line1, req.Line2, req.City, req.State, req.PostCode, req.Country, req.Phone, time.Now(),
	).Scan(&addr.ID, &addr.UserID, &addr.Label, &addr.Line1, &addr.Line2, &addr.City, &addr.State, &addr.PostCode, &addr.Country, &addr.Phone, &addr.CreatedAt, &addr.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

func (s *AddressService) Delete(ctx context.Context, userID, addressID string) error {
	result, err := s.db.Exec(ctx,
		`DELETE FROM addresses WHERE id = $1 AND user_id = $2`,
		addressID, userID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrAddressNotFound
	}
	return nil
}
