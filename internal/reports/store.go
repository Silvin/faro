package reports

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type store struct {
	pool *pgxpool.Pool
}

func newStore(pool *pgxpool.Pool) *store {
	return &store{pool: pool}
}

// salesReport agrega las ventas del negocio en [from, to). tzMinutes es el offset
// del cliente (Date.getTimezoneOffset()) para calcular la hora local.
func (s *store) salesReport(ctx context.Context, tenantID string, from, to time.Time, tzMinutes int) (SalesReport, error) {
	rep := SalesReport{
		ByPaymentMethod: []PaymentBreakdown{},
		ByCategory:      []CategoryBreakdown{},
		ByHour:          []HourBreakdown{},
	}

	// Resumen.
	if err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*), COALESCE(SUM(total_cents), 0)
		   FROM sales WHERE tenant_id = $1 AND created_at >= $2 AND created_at < $3`,
		tenantID, from, to).Scan(&rep.SalesCount, &rep.TotalCents); err != nil {
		return SalesReport{}, err
	}

	// Por forma de pago.
	pmRows, err := s.pool.Query(ctx,
		`SELECT payment_method, COUNT(*), COALESCE(SUM(total_cents), 0)
		   FROM sales WHERE tenant_id = $1 AND created_at >= $2 AND created_at < $3
		  GROUP BY payment_method ORDER BY payment_method`, tenantID, from, to)
	if err != nil {
		return SalesReport{}, err
	}
	for pmRows.Next() {
		var b PaymentBreakdown
		if err := pmRows.Scan(&b.Method, &b.Count, &b.TotalCents); err != nil {
			pmRows.Close()
			return SalesReport{}, err
		}
		rep.ByPaymentMethod = append(rep.ByPaymentMethod, b)
	}
	pmRows.Close()
	if err := pmRows.Err(); err != nil {
		return SalesReport{}, err
	}

	// Por categoría (line items -> producto -> categoría).
	catRows, err := s.pool.Query(ctx,
		`SELECT COALESCE(c.name, 'Sin categoría'), COALESCE(SUM(si.quantity), 0), COALESCE(SUM(si.line_total_cents), 0)
		   FROM sale_items si
		   JOIN sales s ON s.id = si.sale_id
		   LEFT JOIN products p ON p.id = si.product_id
		   LEFT JOIN categories c ON c.id = p.category_id
		  WHERE s.tenant_id = $1 AND s.created_at >= $2 AND s.created_at < $3
		  GROUP BY COALESCE(c.name, 'Sin categoría')
		  ORDER BY SUM(si.line_total_cents) DESC`, tenantID, from, to)
	if err != nil {
		return SalesReport{}, err
	}
	for catRows.Next() {
		var b CategoryBreakdown
		if err := catRows.Scan(&b.CategoryName, &b.Quantity, &b.TotalCents); err != nil {
			catRows.Close()
			return SalesReport{}, err
		}
		rep.ByCategory = append(rep.ByCategory, b)
	}
	catRows.Close()
	if err := catRows.Err(); err != nil {
		return SalesReport{}, err
	}

	// Por hora (hora local del cliente).
	hrRows, err := s.pool.Query(ctx,
		`SELECT EXTRACT(HOUR FROM (created_at - make_interval(mins => $4)))::int AS hr,
		        COUNT(*), COALESCE(SUM(total_cents), 0)
		   FROM sales WHERE tenant_id = $1 AND created_at >= $2 AND created_at < $3
		  GROUP BY hr ORDER BY hr`, tenantID, from, to, tzMinutes)
	if err != nil {
		return SalesReport{}, err
	}
	for hrRows.Next() {
		var b HourBreakdown
		if err := hrRows.Scan(&b.Hour, &b.Count, &b.TotalCents); err != nil {
			hrRows.Close()
			return SalesReport{}, err
		}
		rep.ByHour = append(rep.ByHour, b)
	}
	hrRows.Close()
	return rep, hrRows.Err()
}
