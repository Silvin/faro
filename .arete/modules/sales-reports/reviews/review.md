# Review — sales-reports
_Revisor: code-reviewer · Fecha: 2026-06-30 · Veredicto: **approve**_

## Alcance
`internal/reports` (model/store/service/handler) + frontend `faro-ui` (`/reports`, `lib/reports`). Reutiliza el filtro por fecha de `/sales`.

## Correctitud y diseño
- ✅ Solo lectura; todas las agregaciones filtran por `tenant_id` y `created_at ∈ [from,to)`.
- ✅ Por categoría vía `sale_items → products → categories`; `NULL → "Sin categoría"`.
- ✅ Hora local del cliente vía `tz` offset (`make_interval`), sin hardcodear zona.
- ✅ Aislamiento por negocio probado.

## Seguridad
- ✅ Sesión requerida; sin escritura; sin datos de otros negocios.

## Observaciones (no bloqueantes)
- ⚠️ **[menor]** El desglose por categoría usa la categoría **actual** del producto (no un snapshot). Si se recategoriza un producto, los reportes históricos se re-atribuyen. Aceptable para MVP; si se requiere histórico exacto, guardar `category_id`/nombre en `sale_items` (candidato a mejora).
- ℹ️ El límite del listado de `/sales` (500 con rango) no afecta a los reportes (agregan en SQL, sin límite).

## Veredicto
**approve.** Agregaciones correctas, tenant-scoped y con manejo de zona horaria; única observación es deuda menor de snapshot de categoría.
