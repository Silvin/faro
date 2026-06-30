# Project Charter — Faro
_Creado: 2026-06-29 · Owner: Silvio · Estado: active_

## Problema / Oportunidad
Las cafeterías administran productos, inventario, ventas, reportes y lealtad de forma manual o con herramientas fragmentadas. No tienen un sistema integrado que incluya punto de venta (POS) y fidelización de clientes.

## Objetivo
Un sistema integral de administración de cafeterías **+ punto de venta** que sistematice los procesos clave. Meta del MVP: **instalable y funcional en un negocio real en 4 días**.

## Usuarios
Personal de la cafetería: dueños, gerentes, administradores, baristas y cajeros.
> **v1 sin roles:** todos los usuarios pueden ver y hacer todo. Los roles y permisos llegan en una versión posterior.

## Alcance inicial
- **En alcance:** login, administración de categorías, administración de productos, punto de venta (comanda/venta, cobro con cálculo de cambio, impresión de ticket), reportes de ventas, sistema de lealtad (puntos + tarjeta digital wallet + WhatsApp).
- **Fuera de alcance (por ahora):** roles y permisos, pagos con tarjeta/terminal (v1 cobra en efectivo con cálculo de cambio), inventario avanzado, multi-sucursal.

## Restricciones
- **Tiempo:** MVP funcional instalable en **4 días** = **M1–M4** (login, categorías, productos, POS).
- **POS v1 online** con **impresora térmica** (ESC/POS); migración a **offline** en una fase posterior.
- **Post-MVP (día 5+):** M5 reportes y M6 lealtad (esta última con dependencias externas: wallet + WhatsApp).
- **Sin roles** en v1 (todos pueden ver y hacer todo).

## Stack
- [x] Default del harness (ver `STANDARDS.md` §5): Go · React+Next.js · PostgreSQL · Docker+GitHub Actions · Fly.io+Neon

## Directorio del proyecto
`ARETE_PROJECT_DIR = /Users/silvio/Projects/faro`

---
> El charter es estable pero editable. La secuencia de trabajo vive en `roadmap.md`.
