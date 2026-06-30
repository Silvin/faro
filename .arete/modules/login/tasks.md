# Breakdown — login
_Origen: modules/login/prd.md · tech-spec: modules/login/tech-spec.md · Fecha: 2026-06-29_

## Tareas

### Backend (Go)
- [ ] **T1 — Migración tenants + users** · Tamaño: S · Depende de: — · Crítica: sí
  - Tablas, restricciones (CHECK super_admin/tenant, único email citext) e índices.
  - Done: migración aplica y revierte; esquema coincide con tech-spec.
- [ ] **T2 — Seed del super admin global** · S · Depende de: T1 · Crítica: sí
  - Comando bootstrap idempotente desde `FARO_SUPERADMIN_*`.
  - Done: corriendo el comando existe el super admin; re-correr no duplica.
- [ ] **T3 — Auth core (bcrypt + JWT en cookie)** · M · Depende de: T1 · Crítica: sí
  - Hash/verify bcrypt; emisión/validación JWT HS256; cookie httpOnly+Secure+SameSite.
  - Done: unit tests de hash y de emisión/validación de sesión.
- [ ] **T4 — Endpoints /auth (login, logout, me) + middleware de sesión y tenant-scope** · M · Depende de: T3 · Crítica: sí
  - Middleware carga usuario y acota por `tenant_id`; super admin cruza.
  - Done: integration tests login OK/401, me, logout, aislamiento entre tenants.
- [ ] **T5 — Endpoints de provisión (/tenants, /users, GET /users)** · M · Depende de: T4
  - Validación, 409 email_taken, 403 si no super admin en /tenants.
  - Done: integration tests de alta de negocio y de usuario, acotados por tenant.
- [ ] **T6 — Rate limiting en /auth/login** · S · Depende de: T4
  - Límite por IP+email; respuesta 429.
  - Done: test que verifica el bloqueo tras N intentos.

### Frontend (Next.js)
- [ ] **T7 — App shell + design-system v0.1** · M · Depende de: — · Crítica: sí
  - Tokens, Button, Input, Form field, Card; layout autenticado con top bar + Logout.
  - Done: componentes con sus estados; shell reutilizable (cimiento).
- [ ] **T8 — Pantalla de Login (P1)** · M · Depende de: T7, T4
  - Form email+password, estados vacío/carga/error/éxito, Enter envía, a11y AA.
  - Done: login real contra la API; error genérico; tests de componente.
- [ ] **T9 — Manejo de sesión (guard de rutas + /auth/me + logout)** · M · Depende de: T8
  - Rutas protegidas, redirección a login si no hay sesión, logout.
  - Done: navegar protegido solo con sesión; logout vuelve a login.
- [ ] **T10 — Pantallas mínimas de provisión (P2 crear negocio / P3 crear usuario)** · M · Depende de: T7, T5
  - Forms con validación inline y feedback de éxito/error.
  - Done: alta de negocio (super admin) y de usuario (admin) desde la UI.

## Orden sugerido / ruta crítica
T1 → T2 → T3 → T4 → (T5, T6) · T7 → T8 → T9 → (T10)
Crítica para "poder loguearse": T1 → T3 → T4 → T7 → T8 → T9.

## Riesgos
- Aislamiento entre tenants: cubrir con tests explícitos (negocio A no ve B).
- Seed del super admin en prod: validar manejo seguro de las variables de entorno (devops).
