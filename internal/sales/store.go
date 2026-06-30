package sales

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"faro/internal/dberr"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrInvalidProduct      = errors.New("invalid product")
	ErrInsufficientPayment = errors.New("insufficient payment")
)

type store struct {
	pool *pgxpool.Pool
}

func newStore(pool *pgxpool.Pool) *store {
	return &store{pool: pool}
}

type computedLine struct {
	productID string
	name      string
	unitCents int
	quantity  int
	lineCents int
}

// createSale registra la venta en una transacción. El total se calcula con los
// precios de los productos del negocio (no se confía en el cliente).
func (s *store) createSale(ctx context.Context, tenantID string, items []LineInput, amountPaidCents int) (Sale, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Sale{}, err
	}
	defer tx.Rollback(ctx)

	var lines []computedLine
	total := 0
	for _, it := range items {
		var name, status string
		var price int
		err := tx.QueryRow(ctx,
			`SELECT name, price_cents, status FROM products WHERE id = $1 AND tenant_id = $2`,
			it.ProductID, tenantID).Scan(&name, &price, &status)
		switch {
		case errors.Is(err, pgx.ErrNoRows), dberr.IsInvalidText(err):
			return Sale{}, ErrInvalidProduct
		case err != nil:
			return Sale{}, err
		}
		if status != "active" {
			return Sale{}, ErrInvalidProduct
		}
		lt := price * it.Quantity
		lines = append(lines, computedLine{it.ProductID, name, price, it.Quantity, lt})
		total += lt
	}

	if amountPaidCents < total {
		return Sale{}, ErrInsufficientPayment
	}
	change := amountPaidCents - total

	var sale Sale
	if err := tx.QueryRow(ctx,
		`INSERT INTO sales (tenant_id, total_cents, amount_paid_cents, change_cents)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id::text, tenant_id::text, total_cents, amount_paid_cents, change_cents, created_at`,
		tenantID, total, amountPaidCents, change).
		Scan(&sale.ID, &sale.TenantID, &sale.TotalCents, &sale.AmountPaidCents, &sale.ChangeCents, &sale.CreatedAt); err != nil {
		return Sale{}, err
	}

	for _, l := range lines {
		var item SaleItem
		if err := tx.QueryRow(ctx,
			`INSERT INTO sale_items (sale_id, product_id, name, unit_price_cents, quantity, line_total_cents)
			 VALUES ($1, $2, $3, $4, $5, $6)
			 RETURNING id::text, product_id::text, name, unit_price_cents, quantity, line_total_cents`,
			sale.ID, l.productID, l.name, l.unitCents, l.quantity, l.lineCents).
			Scan(&item.ID, &item.ProductID, &item.Name, &item.UnitPriceCents, &item.Quantity, &item.LineTotalCents); err != nil {
			return Sale{}, err
		}
		sale.Items = append(sale.Items, item)
	}

	if err := tx.Commit(ctx); err != nil {
		return Sale{}, err
	}
	return sale, nil
}

func (s *store) listByTenant(ctx context.Context, tenantID string) ([]Sale, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id::text, tenant_id::text, total_cents, amount_paid_cents, change_cents, created_at
		   FROM sales WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT 50`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Sale
	for rows.Next() {
		var sale Sale
		if err := rows.Scan(&sale.ID, &sale.TenantID, &sale.TotalCents, &sale.AmountPaidCents, &sale.ChangeCents, &sale.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, sale)
	}
	return out, rows.Err()
}

func (s *store) get(ctx context.Context, tenantID, id string) (Sale, error) {
	var sale Sale
	err := s.pool.QueryRow(ctx,
		`SELECT id::text, tenant_id::text, total_cents, amount_paid_cents, change_cents, created_at
		   FROM sales WHERE id = $1 AND tenant_id = $2`, id, tenantID).
		Scan(&sale.ID, &sale.TenantID, &sale.TotalCents, &sale.AmountPaidCents, &sale.ChangeCents, &sale.CreatedAt)
	switch {
	case errors.Is(err, pgx.ErrNoRows), dberr.IsInvalidText(err):
		return Sale{}, ErrNotFound
	case err != nil:
		return Sale{}, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT id::text, product_id::text, name, unit_price_cents, quantity, line_total_cents
		   FROM sale_items WHERE sale_id = $1 ORDER BY id`, sale.ID)
	if err != nil {
		return Sale{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var item SaleItem
		if err := rows.Scan(&item.ID, &item.ProductID, &item.Name, &item.UnitPriceCents, &item.Quantity, &item.LineTotalCents); err != nil {
			return Sale{}, err
		}
		sale.Items = append(sale.Items, item)
	}
	return sale, rows.Err()
}
