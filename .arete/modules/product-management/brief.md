# Módulo: product-management — Brief
_Creado: 2026-06-29 · Estado: idea · Depende de: category-management_

## Qué resuelve
Administrar el catálogo de productos que se venden en la cafetería.

## Resultado esperado
El personal crea, edita, lista y desactiva productos con su precio, categoría y estado.

## Alcance
- **En:** CRUD de productos (nombre, precio, categoría, estado, imagen opcional).
- **Fuera:** variantes/modificadores (ej: tamaño, leche) — a evaluar; inventario/stock avanzado.

## Reutiliza de foundations
- design-system y arquitectura/auth existentes; relación con categorías (M2).

## Notas / dependencias
- Los productos alimentan el punto de venta (M4) y los reportes (M5).
- Pendiente decidir si "stock/inventario" entra aquí o como módulo aparte.
