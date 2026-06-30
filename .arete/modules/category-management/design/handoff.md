# Handoff — category-management
_Fecha: 2026-06-30 · Reutiliza: foundations/design-system.md (BrightPOS) + app shell de login_

## Pantalla: Categorías (`/categories`, dentro del shell autenticado)
- **Layout:** dos zonas (como `/users`): lista de categorías + formulario "Nueva categoría". Responsivo (apila en móvil).
- **Lista:** cada fila = nombre + chip de estado (active/inactive) + acciones (editar, activar/desactivar). Ordenada por `sortOrder` y nombre.
- **Crear:** campo nombre (+ orden opcional) → botón primary lime.
- **Editar:** inline o modal simple para nombre/orden; toggle activar/desactivar.
- **Estados:** vacío ("Sin categorías"), cargando, error (inline), éxito (refresca la lista).
- Componentes existentes: `Card`, `Input`, `FormField`, `Button` (de `faro-ui/components/ui`).

## Navegación
- Nuevo ítem **"Categorías"** en el sidebar (`Sidebar.tsx`), entre Dashboard y Usuarios.

## Notas
- Sin diseño nuevo de sistema: reutiliza tokens y componentes BrightPOS. "Eliminar" = desactivar (no borrado físico).
