# Breakdown — pos
_Origen: prd.md · tech-spec.md · Fecha: 2026-06-30_

## Backend (Go) — repo `faro`
- [ ] **S1 — Migración 0004 sales** · S · tablas sales + sale_items (FKs, índices, check qty>0); up/down.
- [ ] **S2 — Módulo `internal/sales`** · L · Crítica: sí
  - model (Sale, SaleItem), store (`createSale` transaccional con total calculado en servidor, `list`, `get`), service (validación), handler (POST/GET/GET{id}) bajo RequireSession; reutiliza `dberr`.
- [ ] **S3 — Tests** · M · unit (validación: items vacíos, qty≤0, pago insuficiente) + integración (venta con total correcto, cambio, pago insuficiente rechazado, producto inactivo/ajeno rechazado, aislamiento, ticket via get).

## Frontend (Next.js) — repo `faro-ui`
- [ ] **S4 — Página /pos** · L · grid de productos activos → carrito (cantidades) → total → monto recibido → cambio → "Cobrar" (POST /sales) → **ticket** imprimible; ítem "Punto de venta" en el sidebar.
- [ ] **S5 — Ventas del día** · S · lista simple de ventas recientes (GET /sales) con su total/hora.

## Orden
S1 → S2 → S3 (backend) · luego S4 → S5 (frontend).
