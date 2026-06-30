# Tech Spec — login
_Fecha: 2026-06-29 · Basada en: foundations/architecture.md, ADR-001, ADR-002 · Módulo: M1_

## Modelo de datos
```
tenants
  id          uuid pk
  name        text not null
  status      text not null default 'active'   -- active | suspended
  created_at  timestamptz not null default now()

users
  id            uuid pk
  tenant_id     uuid null references tenants(id)   -- null SOLO para super admin global
  email         citext not null unique             -- único global
  password_hash text not null                      -- bcrypt
  name          text not null
  is_super_admin boolean not null default false
  status        text not null default 'active'     -- active | disabled
  created_at    timestamptz not null default now()

  CHECK (is_super_admin = true OR tenant_id IS NOT NULL)  -- usuario de negocio siempre tiene tenant
```
- Índices: `users(tenant_id)`, único `users(email)`.
- Sin tabla de sesiones (JWT sin estado, ver ADR-002).

## Contratos de API (REST/JSON, cookie de sesión httpOnly)

### POST /auth/login
- Request: `{ "email": string, "password": string }`
- 200: `Set-Cookie: faro_session=<jwt>; HttpOnly; Secure; SameSite=Lax` + `{ user }`
- 401: `{ "code": "invalid_credentials", "message": "Email o contraseña incorrectos" }` (genérico)
- 429: `{ "code": "rate_limited", "message": "Demasiados intentos" }`

### POST /auth/logout
- 204: limpia la cookie de sesión.

### GET /auth/me
- 200: `{ user }` si hay sesión válida · 401 si no.

### POST /tenants  _(solo super admin global)_
- Request: `{ "name": string, "ownerEmail": string, "ownerPassword": string, "ownerName": string }`
- 201: `{ "tenant": {...}, "owner": { user sin password } }`
- 400: `{ "code": "validation_error", ... }` · 409: `{ "code": "email_taken" }` · 403 si no es super admin.

### POST /users  _(admin de negocio; super admin debe indicar tenantId)_
- Request: `{ "email": string, "password": string, "name": string }` (acota al tenant del autenticado)
- 201: `{ user sin password }` · 409 `email_taken` · 400 validación.

### GET /users  _(acotado al tenant del autenticado; super admin puede filtrar por tenantId)_
- 200: `{ "items": [ user... ] }`

**Forma de `user`:** `{ id, tenantId, email, name, isSuperAdmin, status, createdAt }` (nunca incluye `password_hash`).

## Reglas técnicas transversales
- **Hash:** bcrypt cost 12. Comparación en tiempo constante.
- **Sesión:** JWT HS256 `{ sub: userId, tid: tenantId|null, sa: isSuperAdmin, exp }`, ~8h, renovación deslizante; cookie httpOnly+Secure+SameSite=Lax.
- **Scoping por tenant:** un helper de acceso a datos exige `tenant_id` salvo en operaciones de super admin. Ningún query de negocio sin tenant.
- **Seed super admin:** comando de bootstrap que lee `FARO_SUPERADMIN_EMAIL` / `FARO_SUPERADMIN_PASSWORD` de entorno y crea (idempotente) el super admin global.
- **Rate limiting:** login limitado por IP+email (ej. 5/min), respuesta 429.
- **Validación:** email con formato válido; password mínimo 8 caracteres.

## Criterios de aceptación técnicos
- [ ] Migración crea `tenants` y `users` con restricciones e índices descritos.
- [ ] Login emite cookie de sesión válida; `/auth/me` la reconoce; logout la invalida.
- [ ] Toda consulta de negocio queda acotada por `tenant_id`; super admin puede cruzar.
- [ ] Contraseñas hasheadas con bcrypt; nunca se devuelve `password_hash`.
- [ ] Rate limiting activo en `/auth/login`.
- [ ] Tests unit + integration cubren login OK/fallo, aislamiento entre tenants y provisión.
