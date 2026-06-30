.PHONY: run build test vet up down migrate-up migrate-down

# Backend
run:        ## Levanta la API en local
go run ./cmd/api

build:      ## Compila el backend
go build ./...

test:       ## Tests del backend
go test ./...

vet:        ## go vet
go vet ./...

# Entorno local (Docker)
up:         ## Levanta Postgres + backend
	docker compose up --build

down:       ## Detiene el entorno local
	docker compose down

# Migraciones (requieren DATABASE_URL en el entorno)
migrate-up:
	psql "$(DATABASE_URL)" -f migrations/0001_init.up.sql

migrate-down:
	psql "$(DATABASE_URL)" -f migrations/0001_init.down.sql
