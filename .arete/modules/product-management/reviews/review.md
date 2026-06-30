# Review — product-management
_Revisor: code-reviewer · Fecha: 2026-06-30 · Veredicto: **approve**_

## Alcance
`internal/products` + `internal/dberr` + `migrations/0003` + frontend `faro-ui` (`/products`, `lib/products`).

## Correctitud y diseño
- ✅ **Dinero como entero (centavos)** + `CHECK (price_cents > 0)` — evita errores de float.
- ✅ Tenant-scope en todas las consultas; `update` exige `id AND tenant_id` → cross-tenant = 404 (probado).
- ✅ Validación de **categoría por negocio** antes de asignar (`categoryBelongsToTenant`) → 400 si es de otro negocio.
- ✅ `update` con `COALESCE($::uuid, category_id)`: cambia la categoría si se envía, la mantiene si se omite.
- ✅ Soft-disable; `LEFT JOIN` para `categoryName` en el listado.

## Seguridad
- ✅ Sesión requerida; aislamiento por negocio probado; entrada validada; sin datos sensibles.

## Estándares y pruebas
- ✅ **Deuda de M2 atendida:** se introdujo `internal/dberr` (helpers de unicidad/uuid) y se usa en products.
- ⚠️ **[menor]** `auth` y `categories` aún tienen su propio helper equivalente; pueden adoptar `internal/dberr` en un refactor de limpieza (no bloqueante).
- ✅ `go vet` / `tsc` limpios; consistente con el patrón de módulos.

## Veredicto
**approve.** Sin hallazgos de seguridad. Buena práctica en el manejo de dinero. Deuda menor de unificación de helpers, seguimiento opcional.
