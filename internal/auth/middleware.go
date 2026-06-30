package auth

import (
	"context"
	"net/http"
)

type ctxKey int

const userCtxKey ctxKey = iota

// RequireSession exige una sesión válida: lee la cookie, valida el JWT y carga
// al usuario (que debe seguir activo), dejándolo en el contexto de la petición.
func (svc *Service) RequireSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(sessionCookieName)
		if err != nil {
			unauthorized(w)
			return
		}
		claims, err := svc.tokens.parse(c.Value)
		if err != nil {
			unauthorized(w)
			return
		}
		u, err := svc.store.userByID(r.Context(), claims.UserID)
		if err != nil || u.Status != "active" {
			unauthorized(w)
			return
		}
		ctx := context.WithValue(r.Context(), userCtxKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserFromContext recupera el usuario autenticado puesto por RequireSession.
func UserFromContext(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(userCtxKey).(User)
	return u, ok
}
