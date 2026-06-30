# Breakdown — product-management
_Origen: prd.md · tech-spec.md · Fecha: 2026-06-30_

## Backend (Go) — repo `faro`
- [ ] **P1 — `internal/dberr`** · S · util compartido (unique 23505, uuid 22P02). Refactor menor de la deuda de M2.
- [ ] **P2 — Migración 0003 products** · S · Crítica: sí · tabla products + FKs + UNIQUE + CHECK; up/down.
- [ ] **P3 — Módulo `internal/products`** · M · Depende de: P1,P2 · model/store/service/handler; CRUD tenant-scoped; validación de categoría por negocio; LEFT JOIN para categoryName; POST/GET/PATCH bajo RequireSession.
- [ ] **P4 — Tests** · M · unit (validación) + integración (crear/listar con categoría, 409 duplicado, categoría de otro negocio rechazada, desactivar, aislamiento, 404 cross-tenant).

## Frontend (Next.js) — repo `faro-ui`
- [ ] **P5 — Página /products** · M · listar (nombre, precio, categoría, estado) + crear (con dropdown de categorías) + editar + activar/desactivar; ítem "Productos" en el sidebar.

## Orden
P1 → P2 → P3 → P4  ·  luego P5 (frontend).
