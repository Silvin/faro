# Test Plan — pos
_Autor: qa-engineer · Fecha: 2026-06-30 · Veredicto: **PASS**_

## Niveles
- Unit (Go): validación (items vacíos, cantidad ≤0).
- Integración (Go + Postgres real `faro_test`): venta, totales, aislamiento.
- E2E (curl) contra backend real + frontend `/pos` (ticket imprimible).

## Cobertura de criterios de aceptación (PRD)
| Criterio | Prueba | Resultado |
|----------|--------|-----------|
| Venta con varios productos; total correcto (servidor) | `TestCreateSaleComputesTotalAndChange` + E2E #1 | ✅ |
| Cobro con monto ≥ total → registra y calcula cambio | idem (change) + E2E #1 ($140) | ✅ |
| Monto < total → rechazado | `TestInsufficientPayment` + E2E #2 (400) | ✅ |
| Producto inactivo / de otro negocio → rechazado | `TestInactiveOrForeignProductRejected` | ✅ |
| Venta persistida; ticket consultable | `TestGetAndIsolation` + E2E #3/#4 | ✅ |
| A no ve ventas de B | `TestGetAndIsolation` (get 404 + listado aislado) | ✅ |

## Resultados
- **5/5 tests verdes**; suites M1–M3 intactas.
- **E2E:** total $60 (=$30×2) calculado por el servidor, cambio $140; pago insuficiente 400; listado y ticket (líneas) OK.
- Build/vet/tsc verdes.

## Quality gate
**PASS** — todos los criterios cubiertos; 0 bloqueantes.

## Nota
- Impresión vía navegador (`@media print`); ESC/POS raw queda como seguimiento (decidido en el charter).
