# PRD — product-management
_Autor: project-manager · Fecha: 2026-06-30 · Estado: draft · Módulo: M3_

## Problema
Cada cafetería necesita administrar el catálogo de productos que vende (con precio y categoría) para usarlo en el POS y los reportes.

## Usuarios
Personal del negocio (v1 sin roles: cualquier usuario autenticado del negocio).

## Requisitos
### Funcionales
- **F1** Crear producto: nombre, precio, categoría (opcional).
- **F2** Listar productos del negocio (con su categoría).
- **F3** Editar nombre / precio / categoría.
- **F4** Activar/desactivar (sin borrado físico).
### No funcionales
- Precio en **centavos enteros** (sin floats).
- Nombre **único por negocio**.
- **Aislamiento por negocio**; la categoría asignada debe ser del mismo negocio.

## Alcance
- **En:** CRUD de productos (nombre, precio, categoría, estado).
- **Fuera:** variantes/modificadores (tamaño, leche), **stock/inventario** (módulo posterior), carga de imágenes (a evaluar; el modelo deja espacio).

## Criterios de aceptación
- [ ] Crear producto con nombre y precio (> 0); categoría opcional del negocio.
- [ ] Listar solo los productos del negocio, con el nombre de su categoría.
- [ ] Editar nombre/precio/categoría.
- [ ] Desactivar (status=inactive); no se borra físicamente.
- [ ] Nombre duplicado por negocio → 409.
- [ ] Asignar categoría de OTRO negocio → rechazado (400/404).
- [ ] Negocio A no ve ni edita productos de B.

## Dependencias
- M1 (auth, tenant-scope) y M2 (categorías). Cimientos disponibles.
