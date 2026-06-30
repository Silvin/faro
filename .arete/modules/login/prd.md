# PRD — login
_Autor: project-manager · Fecha: 2026-06-29 · Estado: draft (Fase A) · Módulo: M1_

## Problema
El personal de cada cafetería necesita entrar a Faro de forma segura. Como Faro es **multi-negocio**, cada usuario pertenece a un negocio y solo debe ver/operar los datos de ese negocio. Falta el sistema de autenticación base sobre el que se construye todo lo demás.

## Usuarios y casos de uso
- **Super admin global** (operador de Faro): da de alta negocios (cafeterías) y la cuenta dueño de cada uno.
- **Admin/dueño del negocio:** inicia sesión y crea usuarios (baristas, cajeros) de su negocio.
- **Usuario del negocio** (barista/cajero): inicia sesión con email+password y opera Faro.
> **v1 sin roles:** dentro de un negocio, cualquier usuario autenticado puede ver y hacer todo.

## Objetivos y métricas de éxito
- Un usuario inicia sesión en < 3 s con email+password.
- Aislamiento total entre negocios (un usuario nunca ve datos de otro negocio).
- Base de auth reutilizable por todos los módulos siguientes (M2–M6).

## Requisitos
### Funcionales
- **F1** Super admin global sembrado en la instalación; puede iniciar sesión.
- **F2** Super admin global crea un negocio (tenant) + su cuenta admin/dueño (email + password temporal).
- **F3** Un usuario admin de un negocio crea usuarios de su negocio (email + password).
- **F4** Cualquier usuario inicia sesión con **email + password** y obtiene una sesión.
- **F5** Cierre de sesión (logout).
- **F6** Endpoint de "usuario actual" (sesión vigente) para proteger rutas.
- **F7** Aislamiento por negocio: cada consulta queda acotada al negocio del usuario; el super admin global puede cruzar negocios.

### No funcionales
- Contraseñas **hasheadas** (nunca en claro ni en logs).
- Sesión segura (cookie httpOnly + Secure).
- **Rate limiting** básico en el login (anti fuerza bruta).
- Mensajes de error genéricos en credenciales inválidas (no revelar si el email existe).

## Alcance
### En alcance
- Modelo de datos de negocios (tenants) y usuarios.
- Seed del super admin global.
- Login / logout / sesión con email + password.
- Alta de negocio (super admin) y alta de usuario (admin), mínima.

### Fuera de alcance
- **PIN** y cajas individuales (fase posterior, cuando haya caja por persona).
- **Roles y permisos** (v1: todos pueden todo dentro del negocio).
- Recuperación de contraseña por email, 2FA, SSO (a evaluar después).
- Pantallas de gestión avanzada de usuarios (solo lo mínimo para provisionar).

## Criterios de aceptación
- [ ] El super admin global puede iniciar sesión con email+password.
- [ ] El super admin puede crear un negocio y su cuenta dueño.
- [ ] Un admin de negocio puede crear un usuario de su negocio.
- [ ] Un usuario inicia sesión y recibe una sesión válida; rutas protegidas exigen sesión.
- [ ] Un usuario de negocio A no puede ver ni operar datos del negocio B.
- [ ] Logout invalida la sesión (la cookie deja de dar acceso).
- [ ] Contraseñas almacenadas hasheadas; credenciales inválidas devuelven error genérico.
- [ ] El login está protegido con rate limiting básico.

## Dependencias
- Ninguna previa (es el primer módulo). Crea los cimientos `architecture.md` y `design-system.md`.

## Preguntas abiertas
- ¿Cómo se siembra el super admin global en producción? (vía variables de entorno + comando de bootstrap — se define en tech-spec.)
- Recuperación de contraseña: ¿se necesita ya o post-MVP? (asumido post-MVP.)
