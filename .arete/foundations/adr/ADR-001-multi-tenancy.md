# ADR-001 — Estrategia de multi-tenancy
_Fecha: 2026-06-29 · Estado: aceptado_

## Contexto
Faro es multi-negocio: una sola plataforma da servicio a varias cafeterías, cada una con sus usuarios y datos aislados, bajo un super admin global. Necesitamos un modelo de aislamiento que sea seguro, barato y simple para el MVP (instalable en 4 días) y que escale.

## Opciones consideradas
- **A) DB compartida + esquema compartido + `tenant_id`** — una columna `tenant_id` en cada tabla de negocio; la app acota las consultas.
  - Pros: simple, barato (una sola DB en Neon), fácil de operar y migrar.
  - Cons: el aislamiento depende de la disciplina de la app (mitigable con helper + RLS).
- **B) Esquema por tenant** — un schema Postgres por negocio.
  - Pros: mejor aislamiento lógico. Cons: migraciones y operación más complejas; no escala a miles de tenants.
- **C) Base de datos por tenant** — una DB por negocio.
  - Pros: aislamiento máximo. Cons: costo y operación altísimos; inviable para MVP.

## Decisión
**Opción A:** DB compartida, esquema compartido, columna `tenant_id` en tablas de negocio. La API **obliga** el scope por `tenant_id`; el super admin global (`tenant_id` nulo + `is_super_admin`) cruza negocios de forma explícita.

## Consecuencias
- Positivas: arranque rápido y barato; una sola migración; alineado con Neon.
- Negativas / deuda: el aislamiento vive en la capa de aplicación. **Mitigación:** un helper de acceso a datos que exija `tenant_id` y, como endurecimiento futuro, **Row-Level Security (RLS)** de Postgres (candidato a ADR posterior).
