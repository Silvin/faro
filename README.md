# Faro — Backend (API)

API del sistema de administración de cafeterías + punto de venta (POS), multi-negocio.

> **Repos:** este repo (`faro`) es el **backend**. El **frontend** vive en **`faro-ui`** (Next.js) y consume esta API por HTTP. Ver `ADR-004`.
> La definición del producto (charter, roadmap, specs por módulo y decisiones de arquitectura) vive en [`.arete/`](./.arete/), generada con el harness **arete-os**.

## Stack
Go · PostgreSQL · Docker + GitHub Actions · Fly.io + Neon.
Arquitectura: **monolito modular** (módulos internos `auth`, `products`, `sales`, `loyalty`). Ver `.arete/foundations/architecture.md` y `ADR-004`.

## Estructura
```
faro/                # backend (Go, código en la raíz)
├── cmd/api/         # entrypoint del servidor
├── internal/<modulo>/  # auth, products, sales, loyalty
├── migrations/      # SQL (up/down)
├── docker-compose.yml  # Postgres + API local
└── .arete/          # specs del producto
```

## Correr en local
Requisitos: Go 1.25+, Docker, psql.

```bash
# 1. Levantar Postgres + backend
make up                 # (o: docker compose up --build)

# 2. Aplicar migraciones (con la DB arriba)
export DATABASE_URL=postgres://faro:faro@localhost:5432/faro?sslmode=disable
make migrate-up

# 3. Probar
curl localhost:8080/health   # {"status":"ok"}
curl localhost:8080/ready    # {"status":"ready"} si la DB responde
```

Sin Docker, solo el backend:
```bash
make run     # usa DATABASE_URL del entorno
make test    # tests
```

## Estado
Módulo **login** en construcción (Fase B). Ver `.arete/modules/login/state.md`.
