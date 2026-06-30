# Breakdown — category-management
_Origen: prd.md · tech-spec: tech-spec.md · Fecha: 2026-06-30_

## Backend (Go) — repo `faro`
- [ ] **C1 — Migración 0002 categories** · S · Crítica: sí
  - Tabla `categories` + `UNIQUE(tenant_id, name)` + índice. Up/down.
- [ ] **C2 — Módulo `internal/categories`** · M · Depende de: C1 · Crítica: sí
  - model, store (create/list/get/update, todo tenant-scoped), service (validación), handler+rutas (POST/GET/PATCH) bajo `RequireSession`.
  - Done: build/vet verdes.
- [ ] **C3 — Tests** · M · Depende de: C2
  - Unit (validación) + integración (crear, listar, editar, desactivar, 409 duplicado, **aislamiento entre negocios**, 404 cross-tenant).

## Frontend (Next.js) — repo `faro-ui`
- [ ] **C4 — Página /categories** · M · Depende de: C2
  - Listar (ordenado) + crear + editar nombre/orden + activar/desactivar, con el design-system BrightPOS y `lib/api`.
  - Done: build verde; ítem "Categorías" en el sidebar.

## Orden
C1 → C2 → C3  ·  luego C4 (frontend).
