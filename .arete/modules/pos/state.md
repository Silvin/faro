# Estado — Módulo pos (punto de venta)
_Actualizado: 2026-06-30_

## Etapa actual
👁 review ✅ — código completo (S1–S5) + E2E + **QA PASS + review approve**. Cierra el MVP M1–M4. Pendiente: deploy (Fase C).

## Pipeline (gates)
| # | Etapa | Dueño | Gate | Estado |
|---|-------|-------|------|--------|
| 1 | PRD | project-manager | PRD aprobado | ✅ |
| 2 | Diseño | product-designer | handoff completo | ✅ (reusa shell + design-system) |
| 3 | Tech-spec | tech-lead | contratos definidos | ✅ |
| 4 | Tareas | project-manager | tareas atómicas | ✅ |
| 5 | Build | backend/frontend | tests verdes | ✅ **S1–S5 completos** (backend 5/5 + frontend `/pos` + ticket); E2E validado |
| 6 | QA | qa-engineer | gate PASS | ✅ **PASS** — `qa/test-plan.md` |
| 7 | Review | code-reviewer | approve | ✅ **approve** — `reviews/review.md` |
| 8 | Deploy | devops-engineer | desplegado + observable | ⬜ |

## Ajustes (2026-06-30) — validados E2E
- POS **tablet-friendly**: **tabs de categorías** (clic en categoría → sus productos), **buscador** en vivo y **imágenes** en el catálogo.

## Decisión clave
- **El backend calcula el total** con los precios de SUS productos (no confía en el cliente) → seguridad.
- Snapshot de nombre y precio por línea al momento de la venta (histórico inmutable).
- v1 online + ticket imprimible (HTML); ESC/POS raw queda como seguimiento.

## Reutiliza
- Auth + tenant-scope; productos (M3) para las líneas; design-system BrightPOS.
