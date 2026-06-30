package sales

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrValidation = errors.New("validation")

// LineInput es una línea solicitada por el cliente (producto + cantidad).
type LineInput struct {
	ProductID string
	Quantity  int
}

type Service struct {
	store *store
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{store: newStore(pool)}
}

// Create valida la solicitud y registra la venta (el total lo calcula el store).
func (svc *Service) Create(ctx context.Context, tenantID string, items []LineInput, paymentMethod string, amountPaidCents int, customerID *string) (Sale, error) {
	if len(items) == 0 || amountPaidCents < 0 {
		return Sale{}, ErrValidation
	}
	if paymentMethod != "cash" && paymentMethod != "card" {
		return Sale{}, ErrValidation
	}
	for _, it := range items {
		if strings.TrimSpace(it.ProductID) == "" || it.Quantity <= 0 {
			return Sale{}, ErrValidation
		}
	}
	if customerID != nil && strings.TrimSpace(*customerID) == "" {
		customerID = nil
	}
	return svc.store.createSale(ctx, tenantID, items, paymentMethod, amountPaidCents, customerID)
}

func (svc *Service) List(ctx context.Context, tenantID string) ([]Sale, error) {
	return svc.store.listByTenant(ctx, tenantID)
}

func (svc *Service) Get(ctx context.Context, tenantID, id string) (Sale, error) {
	return svc.store.get(ctx, tenantID, id)
}
