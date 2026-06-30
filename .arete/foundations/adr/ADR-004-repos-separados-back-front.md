# ADR-004 — Repos separados backend/frontend desde el MVP
_Fecha: 2026-06-29 · Estado: aceptado · Reemplaza: ADR-003_

## Contexto
ADR-003 había definido monorepo con plan de split futuro. El humano decide **separar desde el MVP**: backend y frontend como **repositorios distintos**, con despliegue y evolución independientes desde el día 1. La comunicación entre ambos es **HTTP**.

## Decisión
- **`faro`** = repo del **backend** (Go, monolito modular). El código vive en la **raíz** del repo (`cmd/`, `internal/`, `migrations/`). Las specs del producto (`.arete/`) permanecen aquí como repo principal.
- **`faro-ui`** = repo del **frontend** (Next.js).
- **Comunicación:** REST/JSON sobre **HTTP**. El backend habilita **CORS** para el origen del frontend con `AllowCredentials` (la sesión viaja en cookie httpOnly).
- **Cookies cross-origin:** en prod, back y front bajo el **mismo sitio registrable** (ej. `api.faro.app` / `app.faro.app`) para que `SameSite=Lax` funcione; si quedaran en sitios distintos, la cookie debe ser `SameSite=None; Secure`.
- **Contrato de API = frontera.** Ningún repo importa código del otro; el contrato se versiona.

## Consecuencias
- Positivas: despliegues, CI/CD y versionado **independientes**; equipos/cadencias separables; alineado con la intención de producto.
- Negativas / deuda: hay que mantener el **contrato de API** sincronizado entre dos repos (sin commit atómico cruzado); se requiere **CORS** y cuidado con cookies cross-origin; dos pipelines en vez de uno.
- El monolito modular se mantiene dentro de `faro`; un módulo podría extraerse como servicio en el futuro (ADR aparte).
