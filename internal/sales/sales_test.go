package sales

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// --- Unit (sin DB) ---

func TestCreateValidatesInput(t *testing.T) {
	svc := NewService(nil)
	ctx := context.Background()
	if _, err := svc.Create(ctx, "t1", nil, 100); err != ErrValidation {
		t.Fatalf("items vacíos: esperaba ErrValidation, obtuvo %v", err)
	}
	if _, err := svc.Create(ctx, "t1", []LineInput{{ProductID: "p1", Quantity: 0}}, 100); err != ErrValidation {
		t.Fatalf("cantidad 0: esperaba ErrValidation, obtuvo %v", err)
	}
}

// --- Integración (DB real) ---

func testSvc(t *testing.T) (svc *Service, pool *pgxpool.Pool, a, b, prodA, prodInactive, prodB string, priceA int) {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL no definido; se omiten tests de integración")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Skipf("pool de test: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Skipf("DB de test no disponible: %v", err)
	}
	if _, err := pool.Exec(ctx, "TRUNCATE sale_items, sales, products, categories, users, tenants RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	priceA = 4500
	pool.QueryRow(ctx, "INSERT INTO tenants (name) VALUES ('A') RETURNING id::text").Scan(&a)
	pool.QueryRow(ctx, "INSERT INTO tenants (name) VALUES ('B') RETURNING id::text").Scan(&b)
	pool.QueryRow(ctx, "INSERT INTO products (tenant_id, name, price_cents) VALUES ($1,'Latte',$2) RETURNING id::text", a, priceA).Scan(&prodA)
	pool.QueryRow(ctx, "INSERT INTO products (tenant_id, name, price_cents, status) VALUES ($1,'Viejo',1000,'inactive') RETURNING id::text", a).Scan(&prodInactive)
	pool.QueryRow(ctx, "INSERT INTO products (tenant_id, name, price_cents) VALUES ($1,'Otro',2000) RETURNING id::text", b).Scan(&prodB)
	return NewService(pool), pool, a, b, prodA, prodInactive, prodB, priceA
}

func TestCreateSaleComputesTotalAndChange(t *testing.T) {
	svc, pool, a, _, prodA, _, _, price := testSvc(t)
	defer pool.Close()

	sale, err := svc.Create(context.Background(), a, []LineInput{{ProductID: prodA, Quantity: 2}}, 10000)
	if err != nil {
		t.Fatalf("crear venta: %v", err)
	}
	if sale.TotalCents != price*2 {
		t.Fatalf("total esperaba %d, obtuvo %d", price*2, sale.TotalCents)
	}
	if sale.ChangeCents != 10000-price*2 {
		t.Fatalf("cambio esperaba %d, obtuvo %d", 10000-price*2, sale.ChangeCents)
	}
	if len(sale.Items) != 1 || sale.Items[0].UnitPriceCents != price || sale.Items[0].Quantity != 2 {
		t.Fatalf("línea incorrecta: %+v", sale.Items)
	}
}

func TestInsufficientPayment(t *testing.T) {
	svc, pool, a, _, prodA, _, _, price := testSvc(t)
	defer pool.Close()
	// paga menos que el total (price*1).
	if _, err := svc.Create(context.Background(), a, []LineInput{{ProductID: prodA, Quantity: 1}}, price-1); err != ErrInsufficientPayment {
		t.Fatalf("pago insuficiente: esperaba ErrInsufficientPayment, obtuvo %v", err)
	}
}

func TestInactiveOrForeignProductRejected(t *testing.T) {
	svc, pool, a, _, _, prodInactive, prodB, _ := testSvc(t)
	defer pool.Close()
	ctx := context.Background()
	if _, err := svc.Create(ctx, a, []LineInput{{ProductID: prodInactive, Quantity: 1}}, 100000); err != ErrInvalidProduct {
		t.Fatalf("producto inactivo: esperaba ErrInvalidProduct, obtuvo %v", err)
	}
	if _, err := svc.Create(ctx, a, []LineInput{{ProductID: prodB, Quantity: 1}}, 100000); err != ErrInvalidProduct {
		t.Fatalf("producto de otro negocio: esperaba ErrInvalidProduct, obtuvo %v", err)
	}
}

func TestGetAndIsolation(t *testing.T) {
	svc, pool, a, b, prodA, _, _, _ := testSvc(t)
	defer pool.Close()
	ctx := context.Background()

	sale, err := svc.Create(ctx, a, []LineInput{{ProductID: prodA, Quantity: 1}}, 5000)
	if err != nil {
		t.Fatalf("crear: %v", err)
	}

	// Negocio A obtiene su venta con líneas (ticket).
	got, err := svc.Get(ctx, a, sale.ID)
	if err != nil || len(got.Items) != 1 {
		t.Fatalf("get A: err=%v items=%d", err, len(got.Items))
	}
	// Negocio B no puede verla.
	if _, err := svc.Get(ctx, b, sale.ID); err != ErrNotFound {
		t.Fatalf("get cross-tenant: esperaba ErrNotFound, obtuvo %v", err)
	}
	// Listado aislado.
	listA, _ := svc.List(ctx, a)
	listB, _ := svc.List(ctx, b)
	if len(listA) != 1 || len(listB) != 0 {
		t.Fatalf("aislamiento listado: A=%d B=%d", len(listA), len(listB))
	}
}
