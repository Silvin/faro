# Módulo: pos (punto de venta) — Brief
_Creado: 2026-06-29 · Estado: idea · Depende de: product-management_

## Qué resuelve
Registrar y cobrar ventas en el mostrador de la cafetería.

## Resultado esperado
El cajero/barista crea una comanda o venta general, agrega productos, cobra (con cálculo de cambio) e imprime el ticket.

## Alcance
- **En:** crear comanda o venta general, agregar/quitar productos, totalizar, cobro en efectivo con cálculo de cambio, impresión de ticket.
- **Fuera:** pagos con tarjeta/terminal, propinas avanzadas, división de cuenta (a evaluar).

## Reutiliza de foundations
- Catálogo de productos (M3), design-system, arquitectura/auth.

## Notas / dependencias
- **Decidido:** v1 **online** + **impresora térmica ESC/POS**. La operación **offline** queda para una fase posterior (el tech-lead deberá planear su impacto arquitectónico al migrar).
