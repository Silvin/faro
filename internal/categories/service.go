package categories

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrValidation indica datos inválidos.
var ErrValidation = errors.New("validation")

// Service expone la lógica de categorías.
type Service struct {
	store *store
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{store: newStore(pool)}
}

func (svc *Service) Create(ctx context.Context, tenantID, name string, sortOrder int, imageURL *string) (Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Category{}, ErrValidation
	}
	return svc.store.create(ctx, tenantID, name, sortOrder, normalizeURL(imageURL))
}

func (svc *Service) List(ctx context.Context, tenantID string) ([]Category, error) {
	return svc.store.listByTenant(ctx, tenantID)
}

func (svc *Service) Get(ctx context.Context, tenantID, id string) (Category, error) {
	return svc.store.get(ctx, tenantID, id)
}

// normalizeURL convierte cadenas vacías en nil.
func normalizeURL(u *string) *string {
	if u == nil || strings.TrimSpace(*u) == "" {
		return nil
	}
	v := strings.TrimSpace(*u)
	return &v
}

// UpdateInput agrupa los cambios parciales (nil = no cambia).
type UpdateInput struct {
	Name      *string
	Status    *string
	SortOrder *int
	ImageURL  *string
}

func (svc *Service) Update(ctx context.Context, tenantID, id string, in UpdateInput) (Category, error) {
	if in.Name != nil {
		n := strings.TrimSpace(*in.Name)
		if n == "" {
			return Category{}, ErrValidation
		}
		in.Name = &n
	}
	if in.Status != nil && *in.Status != "active" && *in.Status != "inactive" {
		return Category{}, ErrValidation
	}
	return svc.store.update(ctx, tenantID, id, in.Name, in.Status, in.SortOrder, normalizeURL(in.ImageURL))
}
