# Handoff — product-management
_Fecha: 2026-06-30 · Reutiliza: foundations/design-system.md (BrightPOS) + app shell_

## Pantalla: Productos (`/products`, shell autenticado)
- **Layout:** dos zonas (como `/categories`): lista + formulario "Nuevo producto". Responsivo.
- **Lista:** fila = nombre · precio (formateado `$xx.xx`) · categoría · chip de estado · acciones (renombrar/editar, activar/desactivar).
- **Crear:** nombre, precio (en la UI en pesos, se envía en centavos), **dropdown de categorías** (las activas del negocio) → botón primary lime.
- **Estados:** vacío, cargando, error inline, éxito (refresca).
- Componentes: `Card`, `Input`, `FormField`, `Button` + un `<select>` para categoría.

## Navegación
- Ítem **"Productos"** en el sidebar, después de "Categorías".

## Notas
- Precio: UI en pesos con 2 decimales; convertir a/desde **centavos** al hablar con la API.
- "Eliminar" = desactivar. Sin imágenes en esta versión.
