package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound indica que no existe el registro buscado.
var ErrNotFound = errors.New("not found")

type store struct {
	pool *pgxpool.Pool
}

func newStore(pool *pgxpool.Pool) *store {
	return &store{pool: pool}
}

// userByEmail devuelve el usuario y su password_hash (para verificar login).
func (s *store) userByEmail(ctx context.Context, email string) (User, string, error) {
	var u User
	var hash string
	err := s.pool.QueryRow(ctx,
		`SELECT id::text, tenant_id::text, email, password_hash, name, is_super_admin, status, created_at
		   FROM users WHERE email = $1`, email).
		Scan(&u.ID, &u.TenantID, &u.Email, &hash, &u.Name, &u.IsSuperAdmin, &u.Status, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, "", ErrNotFound
	}
	if err != nil {
		return User{}, "", err
	}
	return u, hash, nil
}

// userByID carga un usuario por id (sin exponer el hash).
func (s *store) userByID(ctx context.Context, id string) (User, error) {
	var u User
	err := s.pool.QueryRow(ctx,
		`SELECT id::text, tenant_id::text, email, name, is_super_admin, status, created_at
		   FROM users WHERE id = $1`, id).
		Scan(&u.ID, &u.TenantID, &u.Email, &u.Name, &u.IsSuperAdmin, &u.Status, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, err
	}
	return u, nil
}

// superAdminExists indica si ya existe un super admin global con ese email (para el seed idempotente).
func (s *store) superAdminExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND is_super_admin = true)`, email).
		Scan(&exists)
	return exists, err
}

func (s *store) createSuperAdmin(ctx context.Context, email, name, hash string) (User, error) {
	var u User
	err := s.pool.QueryRow(ctx,
		`INSERT INTO users (tenant_id, email, password_hash, name, is_super_admin)
		 VALUES (NULL, $1, $2, $3, true)
		 RETURNING id::text, tenant_id::text, email, name, is_super_admin, status, created_at`,
		email, hash, name).
		Scan(&u.ID, &u.TenantID, &u.Email, &u.Name, &u.IsSuperAdmin, &u.Status, &u.CreatedAt)
	return u, err
}
