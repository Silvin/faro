# Review — pos
_Revisor: code-reviewer · Fecha: 2026-06-30 · Veredicto: **approve**_

## Alcance
`internal/sales` (model/store/service/handler) + `migrations/0004` + frontend `faro-ui` (`/pos`, `lib/sales`, ticket).

## Correctitud y diseño
- ✅ **El total se calcula en el servidor** con los precios de `products` (no se confía en el monto/precio del cliente) — práctica correcta y segura.
- ✅ Venta + líneas en **una transacción** (atómico).
- ✅ **Snapshot** de `name` y `unit_price_cents` por línea → el histórico no cambia si luego cambia el producto.
- ✅ Tenant-scope en todo; `GET /sales/{id}` exige `id AND tenant_id` (cross-tenant 404, probado).
- ✅ Dinero en centavos enteros; pago insuficiente rechazado server-side.

## Seguridad
- ✅ Sesión requerida; aislamiento por negocio probado; productos validados (activos y del negocio); sin confianza en datos del cliente para montos.

## Estándares y pruebas
- ✅ Reutiliza `internal/dberr`; consistente con el patrón de módulos; `go vet`/`tsc` limpios.
- ⚠️ **[menor]** El ticket imprime con el truco de visibilidad `@media print`; suficiente para MVP, mejorable con una ruta/vista de impresión dedicada.
- ℹ️ No descuenta inventario (módulo de inventario es posterior, por diseño).

## Veredicto
**approve.** Sin hallazgos de seguridad; muy buen manejo del cálculo de totales del lado servidor.
