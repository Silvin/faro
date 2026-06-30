package products

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrValidation      = errors.New("validation")
	ErrInvalidCategory = errors.New("invalid category")
)

type Service struct {
	store *store
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{store: newStore(pool)}
}

// normalizeCategory devuelve la categoría validada (o nil) para el negocio.
func (svc *Service) normalizeCategory(ctx context.Context, tenantID string, categoryID *string) (*string, error) {
	if categoryID == nil || strings.TrimSpace(*categoryID) == "" {
		return nil, nil
	}
	id := strings.TrimSpace(*categoryID)
	ok, err := svc.store.categoryBelongsToTenant(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrInvalidCategory
	}
	return &id, nil
}

type CreateInput struct {
	Name       string
	PriceCents int
	CategoryID *string
	ImageURL   *string
}

func (svc *Service) Create(ctx context.Context, tenantID string, in CreateInput) (Product, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || in.PriceCents <= 0 {
		return Product{}, ErrValidation
	}
	cat, err := svc.normalizeCategory(ctx, tenantID, in.CategoryID)
	if err != nil {
		return Product{}, err
	}
	return svc.store.create(ctx, tenantID, cat, name, in.PriceCents, normalizeURL(in.ImageURL))
}

// normalizeURL convierte cadenas vacías en nil.
func normalizeURL(u *string) *string {
	if u == nil || strings.TrimSpace(*u) == "" {
		return nil
	}
	v := strings.TrimSpace(*u)
	return &v
}

func (svc *Service) List(ctx context.Context, tenantID string) ([]Product, error) {
	return svc.store.listByTenant(ctx, tenantID)
}

func (svc *Service) Get(ctx context.Context, tenantID, id string) (Product, error) {
	return svc.store.get(ctx, tenantID, id)
}

type UpdateInput struct {
	Name       *string
	PriceCents *int
	CategoryID *string // nil = no cambia; valor = asigna (validado por negocio)
	Status     *string
	ImageURL   *string // nil = no cambia; valor = asigna
}

func (svc *Service) Update(ctx context.Context, tenantID, id string, in UpdateInput) (Product, error) {
	if in.Name != nil {
		n := strings.TrimSpace(*in.Name)
		if n == "" {
			return Product{}, ErrValidation
		}
		in.Name = &n
	}
	if in.PriceCents != nil && *in.PriceCents <= 0 {
		return Product{}, ErrValidation
	}
	if in.Status != nil && *in.Status != "active" && *in.Status != "inactive" {
		return Product{}, ErrValidation
	}
	if in.CategoryID != nil {
		cat, err := svc.normalizeCategory(ctx, tenantID, in.CategoryID)
		if err != nil {
			return Product{}, err
		}
		in.CategoryID = cat // categoría validada (o nil si venía vacía)
	}
	return svc.store.update(ctx, tenantID, id, in.Name, in.PriceCents, in.CategoryID, in.Status, normalizeURL(in.ImageURL))
}
