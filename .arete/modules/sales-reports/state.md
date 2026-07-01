# Estado — Módulo sales-reports
_Actualizado: 2026-06-30_

## Etapa actual
🔨 in-dev (Fase B) — reportes de ventas agregados, reutiliza ventas (M4) y categorías (M2).

## Pipeline (gates)
| # | Etapa | Dueño | Gate | Estado |
|---|-------|-------|------|--------|
| 1 | PRD | project-manager | PRD aprobado | ✅ |
| 2 | Diseño | product-designer | handoff completo | ✅ (reusa shell + design-system) |
| 3 | Tech-spec | tech-lead | contratos definidos | ✅ |
| 4 | Tareas | project-manager | tareas atómicas | ✅ |
| 5 | Build | backend/frontend | tests verdes | 🟡 en curso |
| 6 | QA | qa-engineer | gate PASS | ⬜ |
| 7 | Review | code-reviewer | approve | ⬜ |
| 8 | Deploy | devops-engineer | desplegado + observable | ⬜ |

## Reutiliza
- Ventas (M4) y `sale_items`; categorías (M2) para agrupar; rango de fechas ya
  agregado a `/sales` (reutilizado del ajuste "Ventas del día").
