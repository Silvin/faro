# Handoff — sales-reports
_Fecha: 2026-06-30 · Reutiliza: foundations/design-system.md (BrightPOS) + app shell_

## Pantalla: Reportes (`/reports`, shell autenticado)
- **Filtro de rango** arriba: botones rápidos **Hoy / Ayer** + (opcional) selector de fechas personalizado.
- **Resumen** (tarjetas): Total vendido (grande, lime) y Número de ventas.
- **Por forma de pago:** lista/tarjetas Efectivo vs Tarjeta (monto y conteo).
- **Por categoría:** lista con nombre, unidades y monto; **barra** proporcional al monto.
- **Por horario:** barras por hora del día (0–23) con el monto; resalta las horas pico.
- **Estados:** cargando, vacío ("Sin ventas en el rango"), error.
- Componentes: `Card`, y barras simples con `div` de ancho proporcional (sin librería de charts por ahora).

## Navegación
- Ítem **"Reportes"** en el sidebar (para dueño/gerente).

## Notas
- Montos en pesos (centavos→pesos). Hora local del navegador (se envía `tz`).
