# Tech Spec — sales-reports
_Fecha: 2026-06-30 · Basada en: foundations/architecture.md · Módulo: M5_

## Datos
No requiere tablas nuevas: agrega sobre `sales` y `sale_items` (join a `products`/`categories`).

## Contrato de API (requiere sesión; acotado al tenant)
### GET /reports/sales?from=&to=&tz=
- `from`, `to`: RFC3339 (rango `[from, to)`). Si faltan → hoy.
- `tz`: offset en minutos (`Date.getTimezoneOffset()`) para el desglose por hora en hora local.
- 200:
```json
{
  "totalCents": 0, "salesCount": 0,
  "byPaymentMethod": [ { "method":"cash", "count":0, "totalCents":0 } ],
  "byCategory":      [ { "categoryName":"Bebidas", "quantity":0, "totalCents":0 } ],
  "byHour":          [ { "hour":9, "count":0, "totalCents":0 } ]
}
```

## Consultas (tenant-scoped, rango de fechas)
- **Resumen:** `COUNT(*)`, `SUM(total_cents)` de `sales`.
- **Por pago:** `GROUP BY payment_method`.
- **Por categoría:** `sale_items` ⋈ `sales` (filtro tenant+fecha) ⋈ `products` ⋈ `categories`; `SUM(quantity)`, `SUM(line_total_cents)`; null → "Sin categoría".
- **Por hora:** `EXTRACT(HOUR FROM (created_at - make_interval(mins => tz)))` → hora local; `COUNT`, `SUM`.

## Reglas técnicas
- Solo lectura; todas las consultas filtran por `tenant_id` y `created_at` en `[from,to)`.
- Módulo `internal/reports` (store de agregaciones + service + handler).

## Criterios de aceptación técnicos
- [ ] Totales correctos por rango; desgloses por pago/categoría/hora coherentes.
- [ ] Aislamiento por negocio probado.
- [ ] Tests de integración (DB real) verdes.
