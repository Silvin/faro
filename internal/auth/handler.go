package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes devuelve el router del módulo auth (se monta en /auth).
func (svc *Service) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/login", svc.handleLogin)
	r.Post("/logout", svc.handleLogout)
	r.With(svc.RequireSession).Get("/me", svc.handleMe)
	return r
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (svc *Service) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "Email y contraseña son requeridos")
		return
	}

	u, err := svc.authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		// Mensaje genérico: no revelar si el email existe.
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "Email o contraseña incorrectos")
		return
	}

	token, exp, err := svc.tokens.issue(Claims{UserID: u.ID, TenantID: u.TenantID, IsSuperAdmin: u.IsSuperAdmin})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "No se pudo iniciar sesión")
		return
	}

	svc.setSessionCookie(w, token, exp)
	writeJSON(w, http.StatusOK, map[string]any{"user": u})
}

func (svc *Service) handleLogout(w http.ResponseWriter, _ *http.Request) {
	svc.clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (svc *Service) handleMe(w http.ResponseWriter, r *http.Request) {
	u, ok := UserFromContext(r.Context())
	if !ok {
		unauthorized(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": u})
}
