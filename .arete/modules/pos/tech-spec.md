# Tech Spec — pos
_Fecha: 2026-06-30 · Basada en: foundations/architecture.md · Módulo: M4_

## Modelo de datos
```
sales
  id                uuid pk
  tenant_id         uuid not null references tenants(id)
  total_cents       integer not null
  amount_paid_cents integer not null
  change_cents      integer not null
  created_at        timestamptz not null default now()

sale_items
  id               uuid pk
  sale_id          uuid not null references sales(id) on delete cascade
  product_id       uuid references products(id)        -- snapshot abajo por si cambia
  name             text not null                       -- snapshot
  unit_price_cents integer not null                    -- snapshot
  quantity         integer not null check (quantity > 0)
  line_total_cents integer not null
```
- Índices: `sales(tenant_id, created_at)`, `sale_items(sale_id)`.
- Migración: `0004_sales` (up/down).

## Contratos de API (requieren sesión; acotadas al tenant)
### POST /sales
- Request: `{ "items": [ { "productId": uuid, "quantity": int>0 } ], "amountPaidCents": int }`
- El servidor: valida productos (del negocio y activos), calcula `total` con SUS precios,
  exige `amountPaidCents >= total`, calcula `change`, persiste venta + líneas (transacción).
- 201: `{ sale }` (con `items`) · 400 `validation_error` (items vacíos, cantidad ≤0, producto inválido) · 400 `insufficient_payment`

### GET /sales
- 200: `{ "items": [ saleResumen... ] }` — ventas del negocio, recientes primero (sin líneas).

### GET /sales/{id}
- 200: `{ sale }` con `items` (para el ticket) · 404 si no es del negocio.

**Formas:**
- `sale`: `{ id, tenantId, totalCents, amountPaidCents, changeCents, createdAt, items? }`
- `saleItem`: `{ id, productId, name, unitPriceCents, quantity, lineTotalCents }`

## Reglas técnicas
- **Total calculado en servidor** (precio tomado de `products`, no del request).
- Snapshot de `name` y `unit_price_cents` por línea.
- Tenant-scope en todo; `GET /sales/{id}` exige `id AND tenant_id`.
- Transacción para venta + líneas (atómico).
- Reutiliza `internal/dberr`.

## Criterios de aceptación técnicos
- [ ] Migración crea `sales` y `sale_items` con FKs e índices.
- [ ] Total y cambio correctos (servidor); pago insuficiente rechazado.
- [ ] Producto inactivo/ajeno rechazado; aislamiento por negocio probado.
- [ ] Tests unit + integración (DB real) verdes.
