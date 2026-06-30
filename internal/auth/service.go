package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrInvalidCredentials se devuelve siempre que el login falla, sin distinguir
// si el email no existe o la contraseña es incorrecta (no filtrar información).
var ErrInvalidCredentials = errors.New("invalid credentials")

// Service expone la lógica de autenticación del módulo.
type Service struct {
	store        *store
	tokens       *tokenManager
	cookieSecure bool
}

// NewService construye el servicio de auth.
func NewService(pool *pgxpool.Pool, jwtSecret string, sessionTTL time.Duration, cookieSecure bool) *Service {
	return &Service{
		store:        newStore(pool),
		tokens:       newTokenManager(jwtSecret, sessionTTL),
		cookieSecure: cookieSecure,
	}
}

// authenticate valida credenciales y devuelve el usuario si son correctas.
func (svc *Service) authenticate(ctx context.Context, email, password string) (User, error) {
	u, hash, err := svc.store.userByEmail(ctx, email)
	if errors.Is(err, ErrNotFound) {
		return User{}, ErrInvalidCredentials
	}
	if err != nil {
		return User{}, err
	}
	if u.Status != "active" || !checkPassword(hash, password) {
		return User{}, ErrInvalidCredentials
	}
	return u, nil
}

// SeedSuperAdmin crea el super admin global si no existe (idempotente). Devuelve
// true si lo creó. Pensado para correr al arranque desde variables de entorno.
func (svc *Service) SeedSuperAdmin(ctx context.Context, email, password string) (bool, error) {
	exists, err := svc.store.superAdminExists(ctx, email)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}
	hash, err := hashPassword(password)
	if err != nil {
		return false, err
	}
	if _, err := svc.store.createSuperAdmin(ctx, email, "Super Admin", hash); err != nil {
		return false, err
	}
	return true, nil
}
