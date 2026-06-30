# Handoff — pos
_Fecha: 2026-06-30 · Reutiliza: foundations/design-system.md (BrightPOS) + app shell_

## Pantalla: Punto de venta (`/pos`, shell autenticado)
Layout de dos columnas (responsivo: apila en móvil):
- **Izquierda — Catálogo:** grid/lista de productos **activos** (nombre + precio). Tocar un producto lo agrega al carrito (o suma cantidad).
- **Derecha — Carrito/Cobro:**
  - Líneas: nombre · cantidad (− n +) · subtotal · quitar.
  - **Total** destacado (lime).
  - Campo **"Monto recibido"** → muestra **Cambio** calculado en vivo.
  - Botón **Cobrar** (primary lime) → llama `POST /sales`.
  - Tras cobrar: **Ticket** (modal/vista imprimible) con negocio, líneas, total, recibido, cambio, fecha/hora + botón **Imprimir** (`window.print`).
- **Estados:** carrito vacío (deshabilita Cobrar), error (pago insuficiente, producto inválido), éxito (muestra ticket y limpia carrito).

## Ventas del día (sección o pestaña)
- Lista simple de ventas recientes (GET /sales): hora · total. Tocar → ver ticket.

## Navegación
- Ítem **"Punto de venta"** en el sidebar (arriba, es la pantalla principal de operación).

## Notas
- Precios en pesos (UI) ↔ centavos (API).
- El **cambio** se muestra en vivo en el front, pero el oficial lo devuelve el servidor.
- Ticket pensado para ancho de impresora térmica (columna angosta); impresión vía navegador.
