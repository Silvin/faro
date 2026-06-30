package auth

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrEmailTaken     = errors.New("email taken")
	ErrValidation     = errors.New("validation")
	ErrTenantRequired = errors.New("tenant required")
	ErrTenantNotFound = errors.New("tenant not found")
)

// CreateTenantInput son los datos para dar de alta un negocio con su dueño.
type CreateTenantInput struct {
	Name          string
	OwnerEmail    string
	OwnerPassword string
	OwnerName     string
}

// CreateTenantWithOwner valida y crea el negocio + el usuario dueño (transaccional).
func (svc *Service) CreateTenantWithOwner(ctx context.Context, in CreateTenantInput) (Tenant, User, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.OwnerEmail = strings.TrimSpace(in.OwnerEmail)
	in.OwnerName = strings.TrimSpace(in.OwnerName)
	if in.Name == "" || in.OwnerName == "" || !validEmail(in.OwnerEmail) || len(in.OwnerPassword) < minPasswordLen {
		return Tenant{}, User{}, ErrValidation
	}
	hash, err := hashPassword(in.OwnerPassword)
	if err != nil {
		return Tenant{}, User{}, err
	}
	return svc.store.createTenantWithOwner(ctx, in.Name, in.OwnerEmail, in.OwnerName, hash)
}

// CreateUser valida y crea un usuario dentro del tenant indicado.
func (svc *Service) CreateUser(ctx context.Context, tenantID, email, password, name string) (User, error) {
	email = strings.TrimSpace(email)
	name = strings.TrimSpace(name)
	if name == "" || !validEmail(email) || len(password) < minPasswordLen {
		return User{}, ErrValidation
	}
	exists, err := svc.store.tenantExists(ctx, tenantID)
	if err != nil {
		return User{}, err
	}
	if !exists {
		return User{}, ErrTenantNotFound
	}
	hash, err := hashPassword(password)
	if err != nil {
		return User{}, err
	}
	return svc.store.createUser(ctx, tenantID, email, name, hash)
}

// ListUsers devuelve los usuarios de un negocio.
func (svc *Service) ListUsers(ctx context.Context, tenantID string) ([]User, error) {
	return svc.store.listUsersByTenant(ctx, tenantID)
}
