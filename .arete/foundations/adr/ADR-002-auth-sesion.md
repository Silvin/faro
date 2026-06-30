# ADR-002 — Estrategia de autenticación y sesión
_Fecha: 2026-06-29 · Estado: aceptado_

## Contexto
El login v1 usa **email + password** (sin PIN, sin roles). Es una app web (Next.js + Go) y necesitamos sesiones seguras, simples y compatibles con el modelo multi-tenant. El PIN y las cajas individuales son una fase posterior.

## Opciones consideradas
- **A) JWT (HS256) en cookie httpOnly + Secure** — sesión sin estado en servidor.
  - Pros: simple, sin tabla de sesiones, seguro frente a XSS (no accesible por JS), funciona bien con SSR de Next.js.
  - Cons: revocación inmediata requiere lista de revocados o expiración corta (aceptable con expiración de turno).
- **B) Sesiones en servidor (tabla de sessions + cookie de id)** — estado en DB.
  - Pros: revocación inmediata. Cons: más infraestructura/estado; innecesario para el MVP.
- **C) Proveedor externo (Auth0/Clerk)** — identidad gestionada.
  - Pros: rápido. Cons: costo, dependencia externa y menos control del modelo multi-tenant propio.

## Decisión
**Opción A:** JWT HS256 firmado, en cookie **httpOnly + Secure + SameSite=Lax**, expiración ~8 h con renovación deslizante. Contraseñas con **bcrypt** (cost 12). Logout limpia la cookie. Rate limiting básico en el endpoint de login. El secreto de firma vive en gestión de secretos (no en repo).

## Consecuencias
- Positivas: sin estado de sesión, simple y seguro para web; alineado con el MVP.
- Negativas / deuda: sin revocación inmediata (mitigado por expiración de turno). Si más adelante se requiere revocar al instante (o soportar PIN/multi-dispositivo), se evaluará pasar a sesiones en servidor en un ADR posterior.
