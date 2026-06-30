package customers

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"faro/internal/dberr"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrPhoneTaken = errors.New("phone taken")
)

type store struct {
	pool *pgxpool.Pool
}

func newStore(pool *pgxpool.Pool) *store {
	return &store{pool: pool}
}

func (s *store) create(ctx context.Context, tenantID, phone, firstName, lastName string) (Customer, error) {
	var c Customer
	err := s.pool.QueryRow(ctx,
		`INSERT INTO customers (tenant_id, phone, first_name, last_name)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id::text, tenant_id::text, phone, first_name, last_name, created_at`,
		tenantID, phone, firstName, lastName).
		Scan(&c.ID, &c.TenantID, &c.Phone, &c.FirstName, &c.LastName, &c.CreatedAt)
	if dberr.IsUniqueViolation(err) {
		return Customer{}, ErrPhoneTaken
	}
	return c, err
}

func (s *store) findByPhone(ctx context.Context, tenantID, phone string) (Customer, error) {
	var c Customer
	err := s.pool.QueryRow(ctx,
		`SELECT id::text, tenant_id::text, phone, first_name, last_name, created_at
		   FROM customers WHERE tenant_id = $1 AND phone = $2`, tenantID, phone).
		Scan(&c.ID, &c.TenantID, &c.Phone, &c.FirstName, &c.LastName, &c.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Customer{}, ErrNotFound
	}
	return c, err
}
