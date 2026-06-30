# Review — login
_Revisor: code-reviewer · Fecha: 2026-06-30 · Veredicto: **approve**_

## Alcance
Backend `faro` (`internal/auth`, `config`, `server`, `cmd/api`) + frontend `faro-ui` (login, shell, sesión, provisión). Diff de los commits T1–T10.

## Correctitud y diseño
- ✅ Separación por archivos clara (model/password/token/store/service/handler/middleware); contratos coherentes con `tech-spec.md`.
- ✅ Provisión transaccional (negocio + dueño), manejo de unicidad (23505 → 409), tenant-scope forzado en consultas.
- ⚠️ **[menor]** `internal/auth/ratelimit.go:hits` crece sin evicción. **Aceptado como deuda para MVP** (mono-instancia, entradas pequeñas); migrar a TTL/Redis al escalar.

## Seguridad
- 🔴 **[mayor] → RESUELTO** Rate limit solo por IP+email permitía *credential stuffing* (muchos emails desde una IP sin tope). **Fix:** segundo limitador por IP (20/min) en `service.go` + `handler.go`.
- 🟠 **[menor] → RESUELTO** `JWT_SECRET` por defecto inseguro de forma silenciosa. **Fix:** `config.UsesDefaultJWTSecret()` + advertencia al arranque en `cmd/api`. Operativo pendiente: definir secreto fuerte en prod (runbook/devops).
- ✅ bcrypt cost 12; `password_hash` nunca se serializa; error de login genérico (no enumera emails); cookie `httpOnly`+`SameSite=Lax`+`Secure` configurable; entrada validada en frontera; **aislamiento por tenant probado**.
- ℹ️ Logout es *stateless* (el JWT sigue válido hasta expirar) — decisión consciente (ADR-002). Para revocación inmediata se evaluaría sesión en servidor (ADR futuro).

## Estándares y pruebas
- ✅ `go vet` limpio, `tsc --noEmit` sin errores; estilo consistente con el resto del repo.
- ✅ Cobertura: cada criterio de aceptación tiene prueba; **se añadió test de logout** (gap detectado en la review).

## Veredicto
**approve.** Los hallazgos de seguridad relevantes se resolvieron en este cambio; el resto es deuda documentada, no bloqueante para el MVP.
