# Test Plan — login
_Autor: qa-engineer · Fecha: 2026-06-30 · Veredicto: **PASS**_

## Alcance y niveles
- **Unit (Go):** hashing bcrypt, emisión/validación de JWT, rate limiter.
- **Integración (Go + Postgres real):** flujo de auth y provisión vía HTTP (httptest contra DB `faro_test`).
- **E2E manual (curl):** backend real + DB; verificación CORS cross-origin (`:3000` → `:8080`).
- **Build:** backend (`go build`/`go vet`) y frontend (`npm run build`, `tsc --noEmit`).

## Cobertura de criterios de aceptación (del PRD)
| Criterio de aceptación | Prueba | Resultado |
|------------------------|--------|-----------|
| Super admin global inicia sesión | `TestLoginFlow` + E2E #1 | ✅ |
| Super admin crea negocio + dueño | `TestProvisioningAndTenantIsolation` + E2E #3 | ✅ |
| Admin de negocio crea usuario | `TestProvisioningAndTenantIsolation` + E2E #6 | ✅ |
| Rutas protegidas exigen sesión | `TestLoginFlow` (/me 401 sin cookie) | ✅ |
| Logout invalida la sesión | `TestLogoutInvalidatesSession` | ✅ |
| Negocio A no ve datos de B | `TestProvisioningAndTenantIsolation` | ✅ |
| Contraseñas hasheadas; error genérico | `TestHashAndCheckPassword`, `TestLoginInvalidCredentials` + E2E #8 (401) | ✅ |
| Rate limiting en login | `TestRateLimiter*` + E2E #10 (429 al 6º) | ✅ |

## Resultados
- **10/10 tests verdes** (unit + integración con DB real `faro_test`).
- **E2E curl:** 10/10 pasos OK (login→cookie→/me, crear negocio, login dueño, crear/listar usuarios, 401, 403, 429).
- **CORS** cross-origin verificado (`Allow-Credentials` + `Set-Cookie`).
- **Builds:** backend y frontend verdes.
- **Visual:** login + shell responsivo confirmados en navegador y teléfono (humano).

## Quality gate
**PASS** — todos los criterios de aceptación cubiertos por al menos una prueba; 0 bloqueantes abiertos.

## Deuda / notas (no bloqueante)
- Rate limiter en memoria (mono-instancia). Multi-instancia → store compartido (Redis).
- Suite de integración corre en CI/local con `TEST_DATABASE_URL` a una DB migrada aparte.
