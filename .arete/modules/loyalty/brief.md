# Módulo: loyalty (sistema de lealtad) — Brief
_Creado: 2026-06-29 · Estado: idea · Depende de: product-management, pos_

## Qué resuelve
Fidelizar clientes con un sistema de puntos administrado desde Faro.

## Resultado esperado
El personal da de alta a un cliente por su número telefónico, le genera una tarjeta digital QR (compatible con wallets iOS/Android), puede enviarla por WhatsApp, y al escanear el QR puede sumar puntos o redimir un premio.

## Alcance
- **En:** alta de cliente por teléfono, acumulación de puntos, tarjeta digital QR para Apple Wallet y Google Wallet, envío de tarjeta por WhatsApp, escaneo de QR para sumar puntos o redimir premio.
- **Fuera:** niveles/tiers, campañas de marketing, cupones avanzados (a evaluar).

## Reutiliza de foundations
- Ventas/POS (M4) para acumular puntos, design-system, arquitectura/auth.

## Notas / dependencias (⚠️ externas, con lead time)
- **Apple Wallet:** requiere cuenta Apple Developer + certificados (PassKit).
- **Google Wallet:** requiere alta en Google Wallet API.
- **WhatsApp:** requiere WhatsApp Business API o un proveedor (ej. Twilio) — costo y verificación.
- Estas dependencias impactan el timeline; conviene iniciarlas en paralelo cuanto antes.
