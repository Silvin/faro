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
| 5 | Build | backend/frontend | tests verdes | 🟡 en curso — skeleton + T1 (migración) ✅ |
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
- ✅ Baseline frontend (`faro-ui`): Next.js + Tailwind (BrightPOS) + cliente HTTP.
- ⬜ **Siguiente (backend `faro`):** T2 seed super admin → T3 auth core (bcrypt+JWT) → T4 /auth + middleware tenant, con tests.
- ⬜ Frontend (`faro-ui`): T7 shell → T8 login → T9 sesión → T10 provisión.
