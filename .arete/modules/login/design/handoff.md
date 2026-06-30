# Handoff — login
_Fecha: 2026-06-29 · Diseño base: modules/login/design/wireframes.md · Sistema: foundations/design-system.md_

## Por pantalla

### P1 — Login
- Layout: card centrada (max-width ~360px), padding `24`, radio `md`, sombra `sm`; fondo `--color-bg` (`#F4F4F2`), card blanca. Logo Faro arriba.
- Componentes: `Input(email)`, `Input(password)` con toggle, `Button(primary, ancho completo)` "Entrar" (**lime con texto oscuro**).
- Comportamiento: submit por click o Enter; durante la petición → `Button.loading` + inputs disabled.
- Estados/validación: email con formato válido; password no vacío. Error de servidor → inline genérico "Email o contraseña incorrectos".
- Edge: tras N intentos fallidos el backend aplica rate limit → mostrar "Demasiados intentos, espera un momento".

### P0 — App shell (BrightPOS)
- **Top bar** blanca: logo Faro (izq) + acciones globales (búsqueda, notificaciones). **Sidebar izquierda** blanca: navegación agrupada (secciones con header muted), ítems icono+label, **activo = pastilla lime `#C4E456` con texto oscuro**; al fondo usuario actual + `Button(ghost)` "Salir". **Área de contenido** con breadcrumb + título.
- Tokens: acento lime con texto oscuro, fuente Inter/Poppins (ver `design-system.md` v0.3). Responsive: sidebar colapsable.
- Comportamiento: "Salir" llama logout y redirige a P1.

### P2 / P3 — Crear negocio / Crear usuario
- Form simple con los campos del wireframe, `Button(primary)` "Crear", validación inline, éxito → toast + limpiar form.
- P2 visible solo para super admin global; P3 para usuarios del negocio (v1: cualquiera autenticado del negocio).

## Tokens / componentes usados
- App shell (sidebar), Button (primary, ghost), Input (email/password/text), Form field, Card, Sidebar nav item, Toast/inline error — de `design-system.md` v0.2.

## Assets
- Logo Faro (placeholder por ahora; pendiente versión final).

## Notas de implementación
- Pensar primero en tablet (táctil, ≥44px). Foco visible y navegación por teclado obligatorios (a11y AA).
- La sesión va por cookie httpOnly: el frontend NO guarda tokens en localStorage; usa el endpoint `/auth/me` para saber si hay sesión.
