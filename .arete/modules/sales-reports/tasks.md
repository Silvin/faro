# Breakdown — sales-reports
_Origen: prd.md · tech-spec.md · Fecha: 2026-06-30_

## Backend (Go) — repo `faro`
- [ ] **R1 — Módulo `internal/reports`** · M · Crítica: sí
  - model (SalesReport + desgloses), store (4 consultas de agregación tenant-scoped + rango), service, handler `GET /reports/sales?from=&to=&tz=` bajo RequireSession.
- [ ] **R2 — Tests** · M · integración (resumen, por pago, por categoría, por hora, aislamiento) con DB real.

## Frontend (Next.js) — repo `faro-ui`
- [ ] **R3 — Página /reports** · M · selector de rango (Hoy / Ayer / personalizado) → muestra resumen, por forma de pago, por categoría y por hora (barras simples). Ítem "Reportes" en el sidebar.

## Orden
R1 → R2 (backend) · luego R3 (frontend).
