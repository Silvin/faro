# Estado — Módulo login
_Actualizado: 2026-06-29_

## Etapa actual
🔨 in-dev (Fase B) — skeleton del monorepo + migración (T1) listos y compilando (build/vet verdes)

## Pipeline (gates)
| # | Etapa | Dueño | Gate | Estado |
|---|-------|-------|------|--------|
| 1 | PRD | project-manager | PRD aprobado | ✅ (pend. tu OK) |
| 2 | Diseño | product-designer | handoff completo | ✅ |
| 3 | Tech-spec | tech-lead | contratos definidos | ✅ |
| 4 | Tareas | project-manager | tareas atómicas | ✅ |
| 5 | Build | backend/frontend | tests verdes | ✅ **T1–T10 completos** (backend `faro` + frontend `faro-ui`); builds y unit tests verdes |
| 6 | QA | qa-engineer | gate PASS | ⬜ (Fase B) |
| 7 | Review | code-reviewer | approve | ⬜ (Fase B) |
| 8 | Deploy | devops-engineer | desplegado + observable | ⬜ (Fase C) |

_Estados de gate: ⬜ pendiente · 🟡 en curso · ✅ cumplido · 🔴 bloqueado_

## Decisiones tomadas (Fase A)
- Multi-tenant (multi-negocio) con super admin global — ver ADR-001.
- Auth v1: **solo email + password** (sin PIN, sin roles) — ver ADR-002.
- Sesión: JWT en cookie httpOnly; password con bcrypt.
- **Repos separados desde el MVP** (ADR-004, reemplaza ADR-003): `faro` = backend, `faro-ui` = frontend, comunicación HTTP + CORS (cookie con credentials).

## Aportes a foundations
- `architecture.md` (base + multi-tenancy + auth) — NACIÓ con este módulo.
- `design-system.md` v0.3 (paleta **lime BrightPOS** + Inter/Poppins + app shell top bar/sidebar agrupado + componentes) — NACIÓ.
- `adr/ADR-001-multi-tenancy.md`, `adr/ADR-002-auth-sesion.md`.

## Progreso Fase B
- ✅ Skeleton backend (`faro`): server, config, db, /health, /ready, CORS.
- ✅ T1 migración tenants/users (up/down).
- ✅ **T2 seed super admin** (idempotente, desde env, al arranque).
- ✅ **T3 auth core**: bcrypt + JWT en cookie httpOnly (unit tests verdes).
- ✅ **T4 endpoints**: POST /auth/login, POST /auth/logout, GET /auth/me + middleware `RequireSession` (tests de integración escritos; corren con DB real vía TEST_DATABASE_URL).
- ✅ **T5 provisión**: `POST /tenants` (solo super admin) + `POST`/`GET /users` (acotado al negocio); alta transaccional negocio+dueño, 409 email duplicado, aislamiento por tenant.
- ✅ **T6 rate limiting** en `/auth/login` (5/min por IP+email → 429).
- ✅ Baseline frontend (`faro-ui`): Next.js + Tailwind (BrightPOS) + cliente HTTP.
- 🎉 **Backend de login COMPLETO (T1–T6).**
- ✅ **T7 shell** (sidebar + topbar BrightPOS) + design-system en código.
- ✅ **T8 login** (`/login`, email+password, estados 401/429).
- ✅ **T9 sesión** (guard `/auth/me`, logout, redirección).
- ✅ **T10 provisión** (`/users` listar+crear, `/tenants/new` super admin).
- 🎉 **Login COMPLETO en código (T1–T10), backend + frontend, build verificado en ambos repos.**
- ✅ **Validación local (2026-06-30):** suite de integración VERDE contra Postgres real (incl. aislamiento entre negocios); flujo HTTP E2E con curl OK — login→cookie→/me, crear negocio, login dueño, crear usuario, scoping por negocio, 401 (pass mala), 403 (dueño crea negocio), 429 (rate limit); **CORS cross-origin** (`:3000`→`:8080`) con `Allow-Credentials` y `Set-Cookie` verificado.
- ⬜ **Pendiente:** confirmación visual en navegador; gates formales de QA (test-plan) y review; luego Fase C (deploy Fly.io + Neon).

> Verificado local: `go build`, `go vet`, unit tests (bcrypt+JWT) ✅. Integración pendiente de correr en CI/Docker (no había daemon en el entorno de desarrollo).
