# Test Plan — category-management
_Autor: qa-engineer · Fecha: 2026-06-30 · Veredicto: **PASS**_

## Niveles
- Unit (Go): validación (nombre requerido, status válido).
- Integración (Go + Postgres real `faro_test`): CRUD + aislamiento por negocio.
- E2E (curl) contra backend real + frontend `/categories`.

## Cobertura de criterios de aceptación (PRD)
| Criterio | Prueba | Resultado |
|----------|--------|-----------|
| Crear categoría en su negocio | `TestCreateListAndIsolation` + E2E #1 | ✅ |
| Listar solo las del negocio, ordenadas | `TestCreateListAndIsolation` (orden por sort_order) | ✅ |
| Editar nombre / orden | `TestUpdateDeactivateAndCrossTenant` | ✅ |
| Desactivar (sin borrado físico) | idem (status=inactive) + E2E #5 | ✅ |
| Nombre duplicado por negocio → 409 | `TestDuplicateNamePerTenant` + E2E #4 | ✅ |
| A no ve/edita lo de B (404 cross-tenant) | `TestCreateListAndIsolation`, `TestUpdateDeactivateAndCrossTenant` | ✅ |

## Resultados
- **5/5 tests verdes** (unit + integración con DB real).
- **E2E:** crear 201, listado ordenado, duplicado 409, desactivar 200, super admin sin negocio 400.
- Build/vet/tsc verdes; suite de auth intacta.

## Quality gate
**PASS** — todos los criterios cubiertos; 0 bloqueantes.
