# PRD — sales-reports
_Autor: project-manager · Fecha: 2026-06-30 · Estado: draft · Módulo: M5_

## Problema
El dueño/gerente necesita ver cómo van las ventas para tomar decisiones: cuánto se vendió, en qué se vendió y a qué horas.

## Usuarios
Dueño/gerente/admin del negocio.

## Requisitos
### Funcionales
- **F1** Reporte por **rango de fechas** (hoy, ayer, o periodo).
- **F2** **Resumen**: total vendido y número de ventas.
- **F3** Desglose **por forma de pago** (efectivo/tarjeta).
- **F4** Desglose **por categoría** de producto (cuánto y cuántas unidades).
- **F5** Desglose **por horario** (ventas por hora del día).
### No funcionales
- **Aislamiento por negocio**; cálculos en el servidor.
- Hora del día en la **zona horaria del usuario**.

## Alcance
- **En:** reporte agregado de ventas por rango (resumen + por pago + por categoría + por hora).
- **Fuera:** export a Excel/PDF, gráficos avanzados, comparativas entre periodos (a evaluar).

## Criterios de aceptación
- [ ] Elegir un rango (hoy por defecto) y ver el total y el número de ventas.
- [ ] Ver el total por forma de pago.
- [ ] Ver el total y unidades por categoría.
- [ ] Ver las ventas por hora del día.
- [ ] Un negocio no ve datos de otro.

## Dependencias
- M4 (ventas) y M2 (categorías). Datos ya disponibles.
