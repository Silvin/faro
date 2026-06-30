# Wireframes — login
_Fecha: 2026-06-29 · Módulo: M1_

## Flujo de usuario
```
[No autenticado] → Pantalla de Login → (credenciales OK) → App (home)
                                      → (credenciales mal) → error inline, sigue en Login
[Autenticado] → cualquier ruta protegida accesible → Logout → vuelve a Login
Provisión: Super admin → "Crear negocio"   |   Admin → "Crear usuario"
```

## Pantallas

### P1 — Login · Objetivo: autenticar con email + password
- Bloques: logo Faro · card central · campo Email · campo Password (con toggle) · botón "Entrar" (ancho completo) · área de error inline.
- Estados:
  - **vacío:** campos limpios, botón habilitado.
  - **carga:** botón en estado `loading`, campos deshabilitados.
  - **error:** mensaje genérico "Email o contraseña incorrectos" (no revela cuál).
  - **éxito:** redirige al home de la app.
- Interacción: Enter envía el formulario; foco inicial en Email.

### P2 — Crear negocio (solo super admin global) · Objetivo: dar de alta una cafetería + su dueño
- Campos: Nombre del negocio · Email del dueño · Password temporal · Nombre del dueño.
- Estados: vacío / carga / error (email ya existe) / éxito (negocio creado).

### P3 — Crear usuario (admin de negocio) · Objetivo: provisionar staff
- Campos: Nombre · Email · Password.
- Estados: vacío / carga / error (email ya existe) / éxito (usuario creado).

### P0 — App shell (base) · Objetivo: layout autenticado con **menú lateral**
- Bloques: **sidebar izquierdo** (logo Faro arriba, navegación con icono+label e ítem activo resaltado, al fondo usuario actual + **Logout**) + **área de contenido** a la derecha con header de página.
- Estilo: inspirado en la referencia **BrightPOS** (menú lateral; colores y tipografía a confirmar).
- Responsive: sidebar colapsable en pantallas chicas (tablet de mostrador / escritorio).
- Este shell es **cimiento**: lo reusan todos los módulos siguientes.

## Notas
- UI pensada para **tablet de mostrador** (táctil, controles grandes) y también escritorio.
- Sin pantallas de roles ni PIN en v1.
