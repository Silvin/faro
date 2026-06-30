package products

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// --- Unit (sin DB) ---

func TestCreateValidatesNameAndPrice(t *testing.T) {
	svc := NewService(nil)
	if _, err := svc.Create(context.Background(), "t1", CreateInput{Name: "  ", PriceCents: 100}); err != ErrValidation {
		t.Fatalf("nombre vacío: esperaba ErrValidation, obtuvo %v", err)
	}
	if _, err := svc.Create(context.Background(), "t1", CreateInput{Name: "Café", PriceCents: 0}); err != ErrValidation {
		t.Fatalf("precio 0: esperaba ErrValidation, obtuvo %v", err)
	}
}

// --- Integración (DB real) ---

func testSvc(t *testing.T) (*Service, *pgxpool.Pool, string, string, string, string) {
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
	if _, err := pool.Exec(ctx, "TRUNCATE products, categories, users, tenants RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	var a, b, catA, catB string
	pool.QueryRow(ctx, "INSERT INTO tenants (name) VALUES ('A') RETURNING id::text").Scan(&a)
	pool.QueryRow(ctx, "INSERT INTO tenants (name) VALUES ('B') RETURNING id::text").Scan(&b)
	pool.QueryRow(ctx, "INSERT INTO categories (tenant_id, name) VALUES ($1,'Bebidas') RETURNING id::text", a).Scan(&catA)
	pool.QueryRow(ctx, "INSERT INTO categories (tenant_id, name) VALUES ($1,'Otras') RETURNING id::text", b).Scan(&catB)
	return NewService(pool), pool, a, b, catA, catB
}

func TestCreateListWithCategoryAndIsolation(t *testing.T) {
	svc, pool, a, b, catA, _ := testSvc(t)
	defer pool.Close()
	ctx := context.Background()

	if _, err := svc.Create(ctx, a, CreateInput{Name: "Latte", PriceCents: 4500, CategoryID: &catA}); err != nil {
		t.Fatalf("crear: %v", err)
	}
	list, _ := svc.List(ctx, a)
	if len(list) != 1 {
		t.Fatalf("A esperaba 1 producto, obtuvo %d", len(list))
	}
	if list[0].CategoryName == nil || *list[0].CategoryName != "Bebidas" {
		t.Fatalf("categoryName esperaba 'Bebidas', obtuvo %v", list[0].CategoryName)
	}
	if list[0].PriceCents != 4500 {
		t.Fatalf("precio esperaba 4500, obtuvo %d", list[0].PriceCents)
	}
	// Aislamiento: B no ve productos de A.
	listB, _ := svc.List(ctx, b)
	if len(listB) != 0 {
		t.Fatalf("B esperaba 0 (aislado), obtuvo %d", len(listB))
	}
}

func TestCategoryFromOtherTenantRejected(t *testing.T) {
	svc, pool, a, _, _, catB := testSvc(t)
	defer pool.Close()
	// Producto en A con una categoría de B -> rechazado.
	if _, err := svc.Create(context.Background(), a, CreateInput{Name: "X", PriceCents: 100, CategoryID: &catB}); err != ErrInvalidCategory {
		t.Fatalf("categoría de otro negocio: esperaba ErrInvalidCategory, obtuvo %v", err)
	}
}

func TestDuplicateName(t *testing.T) {
	svc, pool, a, _, _, _ := testSvc(t)
	defer pool.Close()
	ctx := context.Background()
	if _, err := svc.Create(ctx, a, CreateInput{Name: "Espresso", PriceCents: 3000}); err != nil {
		t.Fatalf("crear: %v", err)
	}
	if _, err := svc.Create(ctx, a, CreateInput{Name: "Espresso", PriceCents: 3500}); err != ErrNameTaken {
		t.Fatalf("duplicado: esperaba ErrNameTaken, obtuvo %v", err)
	}
}

func TestUpdateDeactivateAndCrossTenant(t *testing.T) {
	svc, pool, a, b, _, _ := testSvc(t)
	defer pool.Close()
	ctx := context.Background()

	p, err := svc.Create(ctx, a, CreateInput{Name: "Cookie", PriceCents: 2000})
	if err != nil {
		t.Fatalf("crear: %v", err)
	}
	newPrice := 2500
	inactive := "inactive"
	upd, err := svc.Update(ctx, a, p.ID, UpdateInput{PriceCents: &newPrice, Status: &inactive})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if upd.PriceCents != 2500 || upd.Status != "inactive" {
		t.Fatalf("update no aplicó: %+v", upd)
	}
	// Cross-tenant: B no puede editar el producto de A.
	if _, err := svc.Update(ctx, b, p.ID, UpdateInput{Status: &inactive}); err != ErrNotFound {
		t.Fatalf("cross-tenant: esperaba ErrNotFound, obtuvo %v", err)
	}
}
