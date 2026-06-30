# PRD — pos (punto de venta)
_Autor: project-manager · Fecha: 2026-06-30 · Estado: draft · Módulo: M4_

## Problema
El personal necesita registrar y cobrar ventas en el mostrador: armar una venta, agregar productos, cobrar en efectivo (con cálculo de cambio) e imprimir el ticket.

## Usuarios
Cajeros/baristas del negocio (v1 sin roles).

## Requisitos
### Funcionales
- **F1** Armar una venta agregando productos (con cantidad).
- **F2** Ver el total de la venta en tiempo real.
- **F3** Cobrar en efectivo indicando el **monto recibido** → calcular el **cambio**.
- **F4** Registrar la venta (persistente) y mostrar el **ticket** (imprimible).
- **F5** Ver ventas recientes del día.
### No funcionales
- El **total lo calcula el servidor** con los precios de sus productos (no del cliente).
- Cada línea guarda **snapshot** de nombre y precio (histórico inmutable).
- **Aislamiento por negocio**: solo productos y ventas del propio negocio.

## Alcance
- **En:** crear venta con líneas, cobro efectivo + cambio, ticket imprimible, lista de ventas del día.
- **Fuera:** pago con tarjeta/terminal, propinas, división de cuenta, **offline** (fase posterior), impresión ESC/POS raw (se usa impresión del navegador).

## Criterios de aceptación
- [ ] Crear una venta con varios productos y cantidades; el total es correcto (servidor).
- [ ] Cobrar: con monto recibido ≥ total, se registra la venta y se calcula el cambio.
- [ ] Monto recibido < total → rechazado (no se registra).
- [ ] Producto inactivo o de otro negocio → rechazado.
- [ ] La venta queda persistida y se puede ver su ticket.
- [ ] Negocio A no ve ventas de B.

## Dependencias
- M1 (auth), M3 (productos). Cimientos disponibles.
