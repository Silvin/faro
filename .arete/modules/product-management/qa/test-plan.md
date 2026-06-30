# Test Plan — product-management
_Autor: qa-engineer · Fecha: 2026-06-30 · Veredicto: **PASS**_

## Niveles
- Unit (Go): validación (nombre, precio > 0).
- Integración (Go + Postgres real `faro_test`): CRUD + categoría por negocio + aislamiento.
- E2E (curl) contra backend real + frontend `/products`.

## Cobertura de criterios de aceptación (PRD)
| Criterio | Prueba | Resultado |
|----------|--------|-----------|
| Crear producto (nombre, precio>0, categoría opcional) | `TestCreateListWithCategoryAndIsolation` + E2E #1/#2 | ✅ |
| Listar del negocio con nombre de categoría | `TestCreateListWithCategoryAndIsolation` (categoryName) + E2E #3 | ✅ |
| Editar nombre/precio/categoría | `TestUpdateDeactivateAndCrossTenant` | ✅ |
| Desactivar (sin borrado físico) | idem (status=inactive) | ✅ |
| Nombre duplicado por negocio → 409 | `TestDuplicateName` + E2E #6 | ✅ |
| Categoría de OTRO negocio rechazada | `TestCategoryFromOtherTenantRejected` + E2E #5 (400) | ✅ |
| A no ve/edita productos de B | `TestCreateListWithCategoryAndIsolation`, `TestUpdateDeactivateAndCrossTenant` | ✅ |
| Precio inválido (≤0) → 400 | `TestCreateValidatesNameAndPrice` + E2E #4 | ✅ |

## Resultados
- **5/5 tests verdes** (unit + integración); suites de auth y categories intactas.
- **E2E:** crear con/sin categoría (201), listado con categoryName/precio, precio 0→400, categoría ajena→400, duplicado→409.
- Build/vet/tsc verdes.

## Quality gate
**PASS** — todos los criterios cubiertos; 0 bloqueantes.
