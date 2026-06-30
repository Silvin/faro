# ADR-003 — Estructura de repositorio y estrategia de despliegue
_Fecha: 2026-06-29 · Estado: aceptado_

## Contexto
Faro es un SaaS multi-negocio (Go API + Next.js) con un MVP a 4 días y equipo pequeño. Hay que decidir: (1) ¿microservicios o monolito? y (2) ¿monorepo (back+front juntos) o repos separados? Estas decisiones definen la estructura de carpetas, el CI/CD y el modelo de despliegue.

## Opciones consideradas

### Arquitectura del backend
- **A) Monolito modular** — un solo servicio Go con módulos internos (`auth`, `products`, `sales`, `loyalty`).
  - Pros: simple de desarrollar/desplegar/depurar; transacciones locales; ideal para MVP y equipo chico. Pros de separación vía paquetes con límites claros.
  - Cons: todo escala junto (aceptable a esta escala).
- **B) Microservicios** — un servicio por dominio.
  - Pros: escalado y despliegue independientes. Cons: costo operativo alto (varios deploys, red, observabilidad distribuida, consistencia) — **sobreingeniería** para el MVP.

### Estructura de repos
- **C) Monorepo** — `backend/` + `frontend/` en el repo `faro`.
  - Pros: commits atómicos contrato↔consumo; CI/CD único; una fuente de verdad; alineado con el harness.
  - Cons: el repo crece (manejable con estructura clara).
- **D) Polyrepo** — repos separados para back y front.
  - Pros: aislamiento por equipo/release. Cons: cambios de contrato cruzan dos PRs/repos; más fricción para equipo chico.

## Decisión
**Monolito modular (A)** desplegado desde un **monorepo (C)**: un único repo `faro` con `backend/` (Go, módulos internos) y `frontend/` (Next.js), construidos como dos imágenes Docker. Cada módulo del roadmap = un paquete interno del backend, no un servicio.

## Consecuencias
- Positivas: arranque rápido, CI/CD simple, contrato API y su consumo evolucionan juntos.
- Negativas / deuda: escalado acoplado. **Mitigación:** límites de módulo estrictos (paquetes `internal/<modulo>`) que permitan **extraer un microservicio** más adelante si un dominio lo amerita (candidato a ADR futuro). El monolito modular es reversible hacia servicios; el camino inverso es más caro.

## Plan futuro (acordado)
- **Separar `frontend/` (y/o `backend/`) en repos distintos más adelante.** Para que el split sea barato:
  - `backend/` y `frontend/` se mantienen **totalmente desacoplados**: sin imports cruzados; el **único** acoplamiento permitido es el **contrato de API**.
  - Cada uno con su `Dockerfile`, su CI y sus dependencias propias.
  - El split se hará con `git filter-repo` / `git subtree` **conservando la historia**.
  - Hasta entonces, ningún archivo compartido entre back y front en el monorepo.
