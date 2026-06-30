package customers

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrValidation = errors.New("validation")

type Service struct {
	store *store
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{store: newStore(pool)}
}

func (svc *Service) Create(ctx context.Context, tenantID, phone, firstName, lastName string) (Customer, error) {
	phone = strings.TrimSpace(phone)
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	if phone == "" || firstName == "" || lastName == "" {
		return Customer{}, ErrValidation
	}
	return svc.store.create(ctx, tenantID, phone, firstName, lastName)
}

func (svc *Service) FindByPhone(ctx context.Context, tenantID, phone string) (Customer, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return Customer{}, ErrValidation
	}
	return svc.store.findByPhone(ctx, tenantID, phone)
}
