# Módulo: login — Brief
_Creado: 2026-06-29 · Estado: backlog · Depende de: —_

## Qué resuelve
Permitir que el personal de la cafetería entre al sistema de forma segura.

## Resultado esperado
Un usuario inicia sesión y accede a Faro. (v1 sin roles: todos ven y hacen todo.)

## Alcance
- **En:** registro/alta de usuarios, login con credenciales, sesión segura, cierre de sesión.
- **Fuera:** roles y permisos (versión posterior), recuperación de contraseña avanzada (a evaluar).

## Reutiliza de foundations
- Nace con este módulo: base de `architecture.md` (auth) y `design-system.md`.

## Notas / dependencias
- Es el primer módulo: crea los cimientos de autenticación que reusarán los demás.
