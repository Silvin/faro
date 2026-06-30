# Tech Spec — product-management
_Fecha: 2026-06-30 · Basada en: foundations/architecture.md · Módulo: M3_

## Modelo de datos
```
products
  id          uuid pk
  tenant_id   uuid not null references tenants(id)
  category_id uuid null references categories(id)   -- opcional; del mismo negocio
  name        text not null
  price_cents integer not null                       -- dinero como entero (centavos)
  status      text not null default 'active'         -- active | inactive
  created_at  timestamptz not null default now()

  UNIQUE (tenant_id, name)
  CHECK (price_cents > 0)
```
- Índices: `products(tenant_id)`, `products(category_id)`.
- Migración: `0003_products` (up/down).

## Contratos de API (requieren sesión; acotadas al tenant del usuario)
### POST /products
- Request: `{ "name": string, "priceCents": int, "categoryId"?: uuid|null }`
- 201: `{ product }` · 400 validación (incl. categoría de otro negocio) · 409 `name_taken`

### GET /products
- 200: `{ "items": [ product... ] }` (del negocio; incluye `categoryName`)

### PATCH /products/{id}
- Request: `{ "name"?, "priceCents"?, "categoryId"?, "status"? }`
- 200: `{ product }` · 404 si no es de su negocio · 409 `name_taken` · 400 validación

**Forma de `product`:** `{ id, tenantId, categoryId, categoryName, name, priceCents, status, createdAt }`

## Reglas técnicas
- Tenant-scope en todas las consultas (WHERE tenant_id); update exige `id AND tenant_id` → cross-tenant = 404.
- Si viene `categoryId`, validar que la categoría exista **en el mismo negocio**; si no → 400 `invalid_category`.
- `priceCents` entero > 0; `name` no vacío; `status` ∈ {active, inactive}.
- Listado hace LEFT JOIN a categories para `categoryName`.
- Util compartido `internal/dberr` para detectar unicidad (23505) y uuid inválido (22P02).

## Criterios de aceptación técnicos
- [ ] Migración crea `products` con FKs, UNIQUE(tenant_id,name) y CHECK precio>0.
- [ ] CRUD tenant-scoped; categoría validada por negocio; isolation probado.
- [ ] Tests unit + integración (DB real) verdes.
