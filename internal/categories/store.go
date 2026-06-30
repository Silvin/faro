package categories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
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

func pgCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == code
}

func (s *store) create(ctx context.Context, tenantID, name string, sortOrder int) (Category, error) {
	var c Category
	err := s.pool.QueryRow(ctx,
		`INSERT INTO categories (tenant_id, name, sort_order)
		 VALUES ($1, $2, $3)
		 RETURNING id::text, tenant_id::text, name, status, sort_order, created_at`,
		tenantID, name, sortOrder).
		Scan(&c.ID, &c.TenantID, &c.Name, &c.Status, &c.SortOrder, &c.CreatedAt)
	if pgCode(err, "23505") {
		return Category{}, ErrNameTaken
	}
	return c, err
}

func (s *store) listByTenant(ctx context.Context, tenantID string) ([]Category, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id::text, tenant_id::text, name, status, sort_order, created_at
		   FROM categories WHERE tenant_id = $1 ORDER BY sort_order, name`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.TenantID, &c.Name, &c.Status, &c.SortOrder, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// update aplica cambios parciales SOLO si la categoría pertenece al tenant
// (WHERE id AND tenant_id). Si no coincide -> ErrNotFound (aísla entre negocios).
func (s *store) update(ctx context.Context, tenantID, id string, name, status *string, sortOrder *int) (Category, error) {
	var c Category
	err := s.pool.QueryRow(ctx,
		`UPDATE categories
		    SET name       = COALESCE($3, name),
		        status     = COALESCE($4, status),
		        sort_order = COALESCE($5, sort_order)
		  WHERE id = $1 AND tenant_id = $2
		  RETURNING id::text, tenant_id::text, name, status, sort_order, created_at`,
		id, tenantID, name, status, sortOrder).
		Scan(&c.ID, &c.TenantID, &c.Name, &c.Status, &c.SortOrder, &c.CreatedAt)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return Category{}, ErrNotFound
	case pgCode(err, "22P02"): // id con formato de uuid inválido
		return Category{}, ErrNotFound
	case pgCode(err, "23505"):
		return Category{}, ErrNameTaken
	case err != nil:
		return Category{}, err
	}
	return c, nil
}
