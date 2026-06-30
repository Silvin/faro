# Tech Spec — category-management
_Fecha: 2026-06-30 · Basada en: foundations/architecture.md · Módulo: M2_

## Modelo de datos
```
categories
  id          uuid pk
  tenant_id   uuid not null references tenants(id)
  name        text not null
  status      text not null default 'active'   -- active | inactive
  sort_order  int  not null default 0
  created_at  timestamptz not null default now()

  UNIQUE (tenant_id, name)        -- nombre único por negocio
```
- Índice: `categories(tenant_id)`.
- Migración: `0002_categories` (up/down).

## Contratos de API (todas requieren sesión; acotadas al tenant del usuario)
> El `tenant_id` se toma de la sesión (`auth.UserFromContext`). Un super admin global (sin tenant) recibe 400.

### POST /categories
- Request: `{ "name": string, "sortOrder"?: int }`
- 201: `{ category }` · 400 validación · 409 `name_taken` · 401 sin sesión

### GET /categories
- 200: `{ "items": [ category... ] }` (del negocio, orden `sort_order, name`)

### PATCH /categories/{id}
- Request: `{ "name"?: string, "status"?: "active"|"inactive", "sortOrder"?: int }`
- 200: `{ category }` · 404 si no existe **en su negocio** · 409 `name_taken` · 400 validación

**Forma de `category`:** `{ id, tenantId, name, status, sortOrder, createdAt }`

## Reglas técnicas
- Todas las consultas filtran por `tenant_id` (WHERE tenant_id = sesión). `getByID`/`update` exigen `id AND tenant_id` → un negocio no toca lo de otro (404 si no coincide).
- Validación: `name` no vacío (trim); `status` en {active, inactive}.
- Sin borrado físico: "eliminar" = `status = 'inactive'` vía PATCH.

## Criterios de aceptación técnicos
- [ ] Migración crea `categories` con `UNIQUE(tenant_id, name)` e índice.
- [ ] CRUD acotado por tenant; isolation probado (negocio A no ve/edita B).
- [ ] 409 en nombre duplicado por negocio; 404 al tocar categoría de otro negocio.
- [ ] Tests unit + integración (DB real) verdes.
