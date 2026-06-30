# Design System — Faro
_Versión: 0.3 · Fecha: 2026-06-29 · Nace con: módulo login (M1) · Crece con cada módulo_

> **Referencia visual:** _BrightPOS — Point of Sale Dashboard UI_ (Dribbble shot 27026796).
> Identidad: **acento lime/verde-limón** sobre **fondo claro**, sidebar **blanco** con navegación agrupada, ítem activo como **pastilla lime con texto oscuro**, tipografía sans limpia. Estética luminosa y amigable, pensada para mostrador.
>
> Valores de color aproximados a partir de la imagen de referencia; afinar con el archivo original si se requiere precisión exacta.

## Tokens
### Color
- `--color-accent` (lime, **primario/acción**): `#C4E456`
- `--color-accent-strong` (hover/activo): `#B2D63F`
- `--color-on-accent` (texto/iconos sobre lime): `#1A1A1A`  ← **texto oscuro sobre el lime, no blanco**
- `--color-bg` (fondo de página): `#F4F4F2`
- `--color-surface` (sidebar, cards, top bar): `#FFFFFF`
- `--color-text`: `#1A1A1A` · `--color-text-muted`: `#8C8C8C`
- `--color-border`: `#E6E6E4`
- `--color-danger`: `#C0392B` · `--color-success`: `#2E7D32`

### Tipografía
- Familia base UI: **Inter** (excelente para precios/números del POS).
- Marca/encabezados pueden usar una sans redondeada (**Poppins**) para el toque amigable de la referencia.
- Escala (px): `12 · 14 · 16 · 20 · 24 · 32`. Pesos: 400 / 500 / 600 (logo en 700).

### Espaciado / radios / sombras
- Espaciado (px): `4 · 8 · 12 · 16 · 24 · 32`.
- Radios: `sm 8 · md 12 · lg 16` (la referencia usa esquinas bien redondeadas en pastillas y cards).
- Sombra: `sm` muy sutil en cards; el layout se apoya más en bordes claros que en sombras.

## Accesibilidad (reglas base)
- Contraste mínimo **AA**. ⚠️ El lime es claro: **texto sobre lime siempre oscuro** (`--color-on-accent`), nunca blanco.
- Foco visible en todos los controles. Áreas táctiles ≥ 44×44 px (tablet de mostrador).
- Inputs siempre con `<label>` asociado.

## Layout — App shell (cimiento, estilo BrightPOS)
- **Top bar** (blanco): logo **Faro** (izq) + acciones globales (búsqueda, notificaciones, menú "…").
- **Sidebar izquierdo** (blanco, fijo): navegación **agrupada por secciones** con encabezados en gris muted (ej. "Main Menu", "Support"). Cada ítem = icono outline + label. **Ítem activo = pastilla lime (`--color-accent`) con texto oscuro y radio `lg`.**
- **Área de contenido**: breadcrumb (ej. "Dashboard › New Transaction") + título de página + acciones a la derecha (ej. botón "Back" con borde y icono lime).
- Responsive: sidebar colapsable en pantallas chicas. Este shell lo reusan TODOS los módulos.

## Componentes (los que usa login)
### Button — variantes y estados
- `primary`: **fondo lime + texto oscuro**; hover → `accent-strong`. Estados: default / hover / active / disabled / loading.
- `outline`: fondo blanco, borde claro, **icono/acento lime** (como "Back"/"Notification" de la referencia).
- `ghost`: sin fondo, para acciones secundarias (ej. "Salir").
### Input (text / email / password) — estados: default / focus / error / disabled
- Con label, mensaje de error; password con toggle de visibilidad.
### Form field — label + input + texto de error.
### Card / Surface — blanco, borde claro, radio `md`.
### Sidebar nav item — estados: default / hover / **activo (pastilla lime, texto oscuro)** / disabled.
### Section header (sidebar) — texto muted, mayúscula/espaciado, separa grupos ("Main Menu", "Support").

## Patrones
- **Acción principal = lime con texto oscuro.** Acentos (activos, "+", chips seleccionados) en lime.
- **Formulario:** label arriba, error debajo, acción principal a ancho completo en móvil/tablet.
- **Feedback:** errores como inline/toast; mensajes genéricos en credenciales.

## Versionado
- Estable: paleta BrightPOS (lime), tipografía (Inter/Poppins), tokens, Button (primary/outline/ghost), Input, Form field, Card, App shell (top bar + sidebar agrupado), Sidebar nav item, Section header.
- Próximas extensiones (otros módulos): chips de categoría (estilo tabs lime), stepper de cantidad ("− n +" con + lime), grid de cards de producto, breadcrumb.
