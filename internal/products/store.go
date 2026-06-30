package products

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"faro/internal/dberr"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrNameTaken = errors.New("name taken")
)

type store struct {
	pool *pgxpool.Pool
}

func newStore(pool *pgxpool.Pool) *store {
	return &store{pool: pool}
}

// categoryBelongsToTenant valida que la categoría sea del mismo negocio.
func (s *store) categoryBelongsToTenant(ctx context.Context, tenantID, categoryID string) (bool, error) {
	var ok bool
	err := s.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1 AND tenant_id = $2)`,
		categoryID, tenantID).Scan(&ok)
	if dberr.IsInvalidText(err) { // uuid mal formado -> tratamos como inexistente
		return false, nil
	}
	return ok, err
}

func (s *store) create(ctx context.Context, tenantID string, categoryID *string, name string, priceCents int) (Product, error) {
	var p Product
	err := s.pool.QueryRow(ctx,
		`INSERT INTO products (tenant_id, category_id, name, price_cents)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id::text, tenant_id::text, category_id::text, name, price_cents, status, created_at`,
		tenantID, categoryID, name, priceCents).
		Scan(&p.ID, &p.TenantID, &p.CategoryID, &p.Name, &p.PriceCents, &p.Status, &p.CreatedAt)
	if dberr.IsUniqueViolation(err) {
		return Product{}, ErrNameTaken
	}
	return p, err
}

func (s *store) listByTenant(ctx context.Context, tenantID string) ([]Product, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT p.id::text, p.tenant_id::text, p.category_id::text, c.name,
		        p.name, p.price_cents, p.status, p.created_at
		   FROM products p
		   LEFT JOIN categories c ON c.id = p.category_id
		  WHERE p.tenant_id = $1
		  ORDER BY p.name`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.TenantID, &p.CategoryID, &p.CategoryName,
			&p.Name, &p.PriceCents, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// update aplica cambios parciales solo si el producto es del negocio.
// categoryID: nil = no cambia; valor = asigna (validado antes en el service).
func (s *store) update(ctx context.Context, tenantID, id string, name *string, priceCents *int, categoryID, status *string) (Product, error) {
	var p Product
	err := s.pool.QueryRow(ctx,
		`UPDATE products
		    SET name        = COALESCE($3, name),
		        price_cents = COALESCE($4, price_cents),
		        category_id = COALESCE($5::uuid, category_id),
		        status      = COALESCE($6, status)
		  WHERE id = $1 AND tenant_id = $2
		  RETURNING id::text, tenant_id::text, category_id::text, name, price_cents, status, created_at`,
		id, tenantID, name, priceCents, categoryID, status).
		Scan(&p.ID, &p.TenantID, &p.CategoryID, &p.Name, &p.PriceCents, &p.Status, &p.CreatedAt)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return Product{}, ErrNotFound
	case dberr.IsInvalidText(err):
		return Product{}, ErrNotFound
	case dberr.IsUniqueViolation(err):
		return Product{}, ErrNameTaken
	case err != nil:
		return Product{}, err
	}
	return p, nil
}
