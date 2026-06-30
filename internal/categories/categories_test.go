package categories

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// --- Unit (sin DB): validación previa al store ---

func TestCreateRequiresName(t *testing.T) {
	svc := NewService(nil) // el store no se usa cuando la validación falla
	if _, err := svc.Create(context.Background(), "t1", "   ", 0); err != ErrValidation {
		t.Fatalf("nombre vacío: esperaba ErrValidation, obtuvo %v", err)
	}
}

func TestUpdateRejectsInvalidStatus(t *testing.T) {
	svc := NewService(nil)
	bad := "borrado"
	if _, err := svc.Update(context.Background(), "t1", "id", UpdateInput{Status: &bad}); err != ErrValidation {
		t.Fatalf("status inválido: esperaba ErrValidation, obtuvo %v", err)
	}
}

// --- Integración (DB real) ---

func testSvc(t *testing.T) (*Service, *pgxpool.Pool, string, string) {
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
	if _, err := pool.Exec(ctx, "TRUNCATE categories, users, tenants RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	var a, b string
	if err := pool.QueryRow(ctx, "INSERT INTO tenants (name) VALUES ('Cafe A') RETURNING id::text").Scan(&a); err != nil {
		t.Fatalf("tenant A: %v", err)
	}
	if err := pool.QueryRow(ctx, "INSERT INTO tenants (name) VALUES ('Cafe B') RETURNING id::text").Scan(&b); err != nil {
		t.Fatalf("tenant B: %v", err)
	}
	return NewService(pool), pool, a, b
}

func TestCreateListAndIsolation(t *testing.T) {
	svc, pool, a, b := testSvc(t)
	defer pool.Close()
	ctx := context.Background()

	if _, err := svc.Create(ctx, a, "Bebidas calientes", 1); err != nil {
		t.Fatalf("crear: %v", err)
	}
	if _, err := svc.Create(ctx, a, "Panadería", 0); err != nil {
		t.Fatalf("crear: %v", err)
	}

	// Negocio A ve 2, ordenadas por sort_order (Panadería 0, Bebidas 1).
	listA, _ := svc.List(ctx, a)
	if len(listA) != 2 {
		t.Fatalf("A esperaba 2, obtuvo %d", len(listA))
	}
	if listA[0].Name != "Panadería" {
		t.Fatalf("orden incorrecto: %s primero", listA[0].Name)
	}

	// Negocio B no ve nada (aislamiento).
	listB, _ := svc.List(ctx, b)
	if len(listB) != 0 {
		t.Fatalf("B esperaba 0 (aislado), obtuvo %d", len(listB))
	}
}

func TestDuplicateNamePerTenant(t *testing.T) {
	svc, pool, a, _ := testSvc(t)
	defer pool.Close()
	ctx := context.Background()
	if _, err := svc.Create(ctx, a, "Café", 0); err != nil {
		t.Fatalf("crear: %v", err)
	}
	if _, err := svc.Create(ctx, a, "Café", 0); err != ErrNameTaken {
		t.Fatalf("duplicado: esperaba ErrNameTaken, obtuvo %v", err)
	}
}

func TestUpdateDeactivateAndCrossTenant(t *testing.T) {
	svc, pool, a, b := testSvc(t)
	defer pool.Close()
	ctx := context.Background()

	c, err := svc.Create(ctx, a, "Postres", 0)
	if err != nil {
		t.Fatalf("crear: %v", err)
	}

	// Editar nombre + desactivar.
	newName := "Postres y dulces"
	inactive := "inactive"
	upd, err := svc.Update(ctx, a, c.ID, UpdateInput{Name: &newName, Status: &inactive})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if upd.Name != newName || upd.Status != "inactive" {
		t.Fatalf("update no aplicó: %+v", upd)
	}

	// Cross-tenant: el negocio B no puede tocar la categoría de A -> ErrNotFound.
	if _, err := svc.Update(ctx, b, c.ID, UpdateInput{Name: &newName}); err != ErrNotFound {
		t.Fatalf("cross-tenant: esperaba ErrNotFound, obtuvo %v", err)
	}
}
