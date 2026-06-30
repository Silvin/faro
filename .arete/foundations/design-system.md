# Design System — Faro
_Versión: 0.2 · Fecha: 2026-06-29 · Nace con: módulo login (M1) · Crece con cada módulo_

> **Referencia visual:** dashboard POS con **menú lateral izquierdo**, inspirado en _BrightPOS — Point of Sale Dashboard UI_ (Dribbble shot 27026796). Los **colores y la tipografía exactos** están **pendientes de ajustar a la referencia** (ver nota al final).
>
> Base mínima para arrancar. Se **extiende** módulo a módulo (no se rehace).

## Tokens
### Color
- `--color-primary`: café/ámbar (marca cafetería) — ej. `#6F4E37`
- `--color-primary-contrast`: `#FFFFFF`
- `--color-bg`: `#FAFAF8` · `--color-surface`: `#FFFFFF`
- `--color-text`: `#1F1B16` · `--color-text-muted`: `#6B6259`
- `--color-border`: `#E7E2DA`
- `--color-danger`: `#C0392B` · `--color-success`: `#2E7D32`
- _Acento/primario y colores del **sidebar**: a ajustar a la paleta de la referencia (BrightPOS)._

### Tipografía
- Familia: **sans moderna** (propuesta: Poppins o Inter) — _a confirmar con la fuente de la referencia_. Escala: `12 · 14 · 16 · 20 · 24 · 32`.
- Pesos: regular 400, medium 500, semibold 600.

### Espaciado / radios / sombras
- Espaciado (px): `4 · 8 · 12 · 16 · 24 · 32`.
- Radios: `sm 6 · md 10 · lg 16`. Sombra: `sm` sutil para tarjetas/inputs en foco.

## Accesibilidad (reglas base)
- Contraste mínimo **AA**. Foco visible en todos los controles.
- Áreas táctiles ≥ 44×44 px (pensado para tablet de mostrador).
- Inputs siempre con `<label>` asociado.

## Layout — App shell con menú lateral (cimiento)
- **Sidebar izquierdo** fijo: logo Faro arriba · ítems de navegación (icono + label, ítem **activo resaltado**) · al fondo, usuario actual + **Logout**. Colapsable en pantallas chicas.
- **Área de contenido** a la derecha, con header de página. Este shell lo reusan TODOS los módulos (POS, productos, reportes…).
- Estilo de menú/colores/tipografía: según la referencia BrightPOS (pendiente de ajuste fino).

## Componentes (los que usa login)
### Button — estados: default / hover / active / disabled / loading
- Variantes: `primary` (acción principal), `ghost` (secundaria).
### Input (text / email / password) — estados: default / focus / error / disabled
- Con label, mensaje de error y, en password, toggle de visibilidad.
### Form field — label + input + texto de error.
### Card / Surface — contenedor con borde y radio `md`.
### Sidebar nav item — estados: default / hover / **activo** / disabled (icono + label).
### Toast / Inline error — feedback de error genérico (ej. credenciales inválidas).

## Patrones
- **Formulario:** label arriba, error debajo del campo, acción principal a ancho completo en móvil/tablet.
- **Feedback:** errores de servidor como inline/toast; nunca revelar detalles sensibles.

## Versionado
- Estable: tokens base, Button, Input, Form field, Card, **App shell (sidebar)**, Sidebar nav item.
- Pendiente de ajuste a la referencia: paleta de acento/sidebar y tipografía exacta.

## Nota — cómo fijar colores y tipografía exactos de la referencia
No pude leer la imagen de Dribbble (se carga por JS). Para clavar los valores exactos, comparte una de estas:
- Los **hex** de los colores + el **nombre de la fuente**, o
- Una **captura** del shot guardada en un archivo (ej. `/Users/silvio/Projects/faro/ref-brightpos.png`) y la **leo** para extraer la paleta.
