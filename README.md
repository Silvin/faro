# Faro

Sistema de administración de cafeterías + punto de venta (POS), multi-negocio.

> La definición del producto (charter, roadmap, specs por módulo y decisiones de arquitectura) vive en [`.arete/`](./.arete/), generada con el harness **arete-os**.

## Stack
Go (backend) · React + Next.js (frontend) · PostgreSQL · Docker + GitHub Actions · Fly.io + Neon.
Arquitectura: **monolito modular** en **monorepo** (`backend/` + `frontend/`). Ver `.arete/foundations/architecture.md` y `ADR-003`.

## Estructura
```
faro/
├── backend/    # Go (monolito modular): cmd/api, internal/<modulo>, migrations
├── frontend/   # Next.js (próximo incremento)
├── docker-compose.yml
└── .arete/     # specs del proyecto
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
