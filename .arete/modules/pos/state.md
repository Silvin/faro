# Estado — Módulo pos (punto de venta)
_Actualizado: 2026-06-30_

## Etapa actual
🔨 in-dev (Fase B) — reutiliza auth (tenant-scope), productos (M3) y design-system.

## Pipeline (gates)
| # | Etapa | Dueño | Gate | Estado |
|---|-------|-------|------|--------|
| 1 | PRD | project-manager | PRD aprobado | ✅ |
| 2 | Diseño | product-designer | handoff completo | ✅ (reusa shell + design-system) |
| 3 | Tech-spec | tech-lead | contratos definidos | ✅ |
| 4 | Tareas | project-manager | tareas atómicas | ✅ |
| 5 | Build | backend/frontend | tests verdes | 🟡 backend S1–S3 ✅ (5/5 tests); falta S4–S5 frontend |
| 6 | QA | qa-engineer | gate PASS | ⬜ |
| 7 | Review | code-reviewer | approve | ⬜ |
| 8 | Deploy | devops-engineer | desplegado + observable | ⬜ |

## Decisión clave
- **El backend calcula el total** con los precios de SUS productos (no confía en el cliente) → seguridad.
- Snapshot de nombre y precio por línea al momento de la venta (histórico inmutable).
- v1 online + ticket imprimible (HTML); ESC/POS raw queda como seguimiento.

## Reutiliza
- Auth + tenant-scope; productos (M3) para las líneas; design-system BrightPOS.
