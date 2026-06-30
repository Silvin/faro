package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// isUniqueViolation detecta el error de Postgres por violación de unicidad (23505).
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func (s *store) tenantExists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tenants WHERE id = $1)`, id).Scan(&exists)
	return exists, err
}

// createTenantWithOwner crea el negocio y su dueño en una transacción.
func (s *store) createTenantWithOwner(ctx context.Context, name, ownerEmail, ownerName, ownerHash string) (Tenant, User, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Tenant{}, User{}, err
	}
	defer tx.Rollback(ctx)

	var t Tenant
	if err := tx.QueryRow(ctx,
		`INSERT INTO tenants (name) VALUES ($1)
		 RETURNING id::text, name, status, created_at`, name).
		Scan(&t.ID, &t.Name, &t.Status, &t.CreatedAt); err != nil {
		return Tenant{}, User{}, err
	}

	var u User
	if err := tx.QueryRow(ctx,
		`INSERT INTO users (tenant_id, email, password_hash, name)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id::text, tenant_id::text, email, name, is_super_admin, status, created_at`,
		t.ID, ownerEmail, ownerHash, ownerName).
		Scan(&u.ID, &u.TenantID, &u.Email, &u.Name, &u.IsSuperAdmin, &u.Status, &u.CreatedAt); err != nil {
		if isUniqueViolation(err) {
			return Tenant{}, User{}, ErrEmailTaken
		}
		return Tenant{}, User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Tenant{}, User{}, err
	}
	return t, u, nil
}

func (s *store) createUser(ctx context.Context, tenantID, email, name, hash string) (User, error) {
	var u User
	err := s.pool.QueryRow(ctx,
		`INSERT INTO users (tenant_id, email, password_hash, name)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id::text, tenant_id::text, email, name, is_super_admin, status, created_at`,
		tenantID, email, hash, name).
		Scan(&u.ID, &u.TenantID, &u.Email, &u.Name, &u.IsSuperAdmin, &u.Status, &u.CreatedAt)
	if isUniqueViolation(err) {
		return User{}, ErrEmailTaken
	}
	return u, err
}

// listUsersByTenant lista usuarios acotados a un negocio (aislamiento por tenant).
func (s *store) listUsersByTenant(ctx context.Context, tenantID string) ([]User, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id::text, tenant_id::text, email, name, is_super_admin, status, created_at
		   FROM users WHERE tenant_id = $1 ORDER BY created_at`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.TenantID, &u.Email, &u.Name, &u.IsSuperAdmin, &u.Status, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
