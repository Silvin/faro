# Arquitectura — Faro
_Fecha: 2026-06-29 · Nace con: módulo login (M1) · Se extiende con cada módulo_

## Visión general
Aplicación web **multi-negocio (multi-tenant)** para administración de cafeterías + POS.

```
[ Next.js (admin + POS web) ]  ──HTTPS/JSON──►  [ Go API ]  ──►  [ PostgreSQL ]
        TypeScript, Tailwind                     net/http+Chi, pgx        (Neon en prod)
```

## Vista de componentes
| Componente | Tecnología | Responsabilidad | Dueño |
|------------|-----------|-----------------|-------|
| Web app | React + Next.js (App Router, TS, Tailwind) | UI de administración y POS; consume la API | frontend-engineer |
| API | Go (Chi router, pgx) | Lógica de negocio, auth, contratos REST | backend-engineer |
| Base de datos | PostgreSQL (Neon en prod) | Persistencia con aislamiento por negocio | backend-engineer |
| Plataforma | Docker + GitHub Actions + Fly.io | Build, deploy, runtime | devops-engineer |

## Estructura de repo y despliegue (decisión base — ver ADR-003)
- **Monolito modular** en Go: un solo servicio con módulos internos (`auth`, `products`, `sales`, `loyalty`), **no microservicios**.
- **Monorepo:** un único repo `faro` con `backend/` y `frontend/`, desplegados como dos imágenes Docker.

```
faro/
├── backend/                 # Go (monolito modular)
│   ├── cmd/api/             # entrypoint del servidor
│   ├── internal/<modulo>/   # auth, products, sales, loyalty (límites estrictos)
│   ├── migrations/
│   └── Dockerfile
├── frontend/                # Next.js (App Router, TS, Tailwind)
│   └── Dockerfile
├── docker-compose.yml       # local: backend + frontend + Postgres
├── fly.backend.toml · fly.frontend.toml
└── .arete/                  # specs del harness
```
> Cada módulo es un paquete interno con frontera clara, para poder extraer un microservicio en el futuro si hiciera falta. Reversible hacia servicios; no al revés.
> **Plan futuro:** separar `frontend/` a su propio repo. Por eso `backend/` y `frontend/` se mantienen **desacoplados** (sin imports cruzados; acoplados solo por el contrato de API), para que el split sea barato y conserve historia.

## Multi-tenancy (decisión base — ver ADR-001)
- **Modelo:** base de datos compartida, esquema compartido, columna **`tenant_id`** en las tablas con datos de negocio.
- **Aislamiento:** la API acota **toda** consulta al `tenant_id` del usuario autenticado. El **super admin global** (`tenant_id` nulo, `is_super_admin`) puede operar sobre cualquier negocio de forma explícita.
- **Endurecimiento futuro:** Row-Level Security (RLS) de Postgres como segunda barrera.

## Autenticación (decisión base — ver ADR-002)
- **Login:** email + password. Password hasheado con **bcrypt**.
- **Sesión:** **JWT (HS256)** firmado, guardado en cookie **httpOnly + Secure + SameSite=Lax**; expiración ~8 h (un turno), renovación deslizante.
- **Middleware:** valida la sesión, carga el usuario y su `tenant_id`, y acota el scope. (PIN y cajas: fase posterior.)

## Límites y contratos
- **Web ↔ API:** REST/JSON. La cookie de sesión viaja en cada request (no hay tokens en localStorage).
- **API ↔ DB:** acceso solo desde la API; el frontend nunca habla con la DB.
- Contrato detallado por módulo en `modules/<m>/tech-spec.md`.

## Requisitos no funcionales y cómo se cumplen
- **Seguridad:** hashing de contraseñas, cookie httpOnly/Secure, rate limiting en login, aislamiento por tenant.
- **Rendimiento:** índices por `tenant_id` y por claves de búsqueda; consultas acotadas.
- **Despliegue:** contenedores Docker en Fly.io; Postgres gestionado en Neon (ver `runbook.md`).

## Riesgos
- Fuga entre tenants si una consulta olvida el `tenant_id` → mitigar con un helper de acceso a datos que **exija** tenant scope + (futuro) RLS.
- Secreto de firma del JWT debe estar en gestión de secretos (no en repo).
