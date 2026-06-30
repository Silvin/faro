# Review — category-management
_Revisor: code-reviewer · Fecha: 2026-06-30 · Veredicto: **approve**_

## Alcance
`internal/categories` (model/store/service/handler) + `migrations/0002` + frontend `faro-ui` (`/categories`, `lib/categories`).

## Correctitud y diseño
- ✅ CRUD acotado por negocio; `update`/`get` exigen `id AND tenant_id` → cross-tenant devuelve 404 (probado).
- ✅ Soft-disable (status) en vez de borrado físico, alineado con el PRD.
- ✅ `sort_order` como `*int` en update distingue "no cambiar" de "poner 0".
- ✅ Manejo de `22P02` (uuid inválido) → 404 (no 500).

## Seguridad
- ✅ Todas las consultas filtran por `tenant_id`; sesión requerida (`RequireSession`); super admin global sin negocio → 400.
- ✅ Sin datos sensibles; entrada validada (nombre no vacío, status en {active,inactive}).

## Estándares y pruebas
- ✅ Consistente con el patrón de `internal/auth`; `go vet` y `tsc` limpios.
- ⚠️ **[menor]** Helper de violación de unicidad duplicado (`pgCode` en categories ≈ `isUniqueViolation` en auth). **Deuda:** extraer a un paquete util compartido cuando haya un 3er módulo (M3).

## Veredicto
**approve.** Sin hallazgos de seguridad; única observación es deuda menor de duplicación, candidata a refactor en M3.
