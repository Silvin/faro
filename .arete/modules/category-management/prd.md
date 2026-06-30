# PRD — category-management
_Autor: project-manager · Fecha: 2026-06-30 · Estado: draft · Módulo: M2_

## Problema
Cada cafetería necesita organizar sus productos en categorías (bebidas calientes, frías, panadería…) para el catálogo, el POS y los reportes.

## Usuarios y casos de uso
- Personal del negocio (dueño/admin/barista) crea y administra las categorías de **su** negocio.
> v1 sin roles: cualquier usuario autenticado del negocio puede administrarlas.

## Objetivos
- CRUD de categorías acotado por negocio, reutilizando la sesión y el aislamiento de login.

## Requisitos
### Funcionales
- **F1** Crear categoría (nombre, orden opcional) en el negocio del usuario.
- **F2** Listar las categorías del negocio (ordenadas por `sortOrder` y nombre).
- **F3** Editar nombre y/o orden de una categoría.
- **F4** Activar/desactivar una categoría (sin borrado físico).
### No funcionales
- Nombre **único por negocio** (no duplicados).
- **Aislamiento por negocio:** un negocio no ve ni edita categorías de otro.

## Alcance
- **En:** CRUD (crear, listar, editar, activar/desactivar) de categorías por negocio.
- **Fuera:** subcategorías/jerarquías, asignación de productos (eso es M3), reordenar drag&drop avanzado.

## Criterios de aceptación
- [ ] Un usuario crea una categoría en su negocio.
- [ ] Lista solo las categorías de su negocio, ordenadas.
- [ ] Edita nombre/orden de una categoría suya.
- [ ] Desactiva (status=inactive) una categoría; no se borra físicamente.
- [ ] Nombre duplicado en el mismo negocio → error (409).
- [ ] Un usuario del negocio A no puede ver ni modificar categorías del negocio B.

## Dependencias
- Login (M1): sesión, `RequireSession`, tenant-scope. Cimientos ya disponibles.
