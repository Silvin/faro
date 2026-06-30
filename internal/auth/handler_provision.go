package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// requireSuperAdmin restringe el acceso al super admin global.
func (svc *Service) requireSuperAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := UserFromContext(r.Context())
		if !ok {
			unauthorized(w)
			return
		}
		if !u.IsSuperAdmin {
			writeError(w, http.StatusForbidden, "forbidden", "Requiere super admin global")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// TenantRoutes: alta de negocios (solo super admin global). Se monta en /tenants.
func (svc *Service) TenantRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(svc.RequireSession)
	r.With(svc.requireSuperAdmin).Post("/", svc.handleCreateTenant)
	return r
}

// UserRoutes: alta y listado de usuarios, acotado al negocio. Se monta en /users.
func (svc *Service) UserRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(svc.RequireSession)
	r.Post("/", svc.handleCreateUser)
	r.Get("/", svc.handleListUsers)
	return r
}

type createTenantRequest struct {
	Name          string `json:"name"`
	OwnerEmail    string `json:"ownerEmail"`
	OwnerPassword string `json:"ownerPassword"`
	OwnerName     string `json:"ownerName"`
}

func (svc *Service) handleCreateTenant(w http.ResponseWriter, r *http.Request) {
	var req createTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", "Cuerpo inválido")
		return
	}
	tenant, owner, err := svc.CreateTenantWithOwner(r.Context(), CreateTenantInput{
		Name: req.Name, OwnerEmail: req.OwnerEmail, OwnerPassword: req.OwnerPassword, OwnerName: req.OwnerName,
	})
	switch {
	case errors.Is(err, ErrValidation):
		writeError(w, http.StatusBadRequest, "validation_error", "Nombre, email del dueño y contraseña (≥ 8) son requeridos")
	case errors.Is(err, ErrEmailTaken):
		writeError(w, http.StatusConflict, "email_taken", "El email ya está registrado")
	case err != nil:
		writeError(w, http.StatusInternalServerError, "internal", "No se pudo crear el negocio")
	default:
		writeJSON(w, http.StatusCreated, map[string]any{"tenant": tenant, "owner": owner})
	}
}

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	TenantID string `json:"tenantId"` // requerido solo si lo crea el super admin global
}

func (svc *Service) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	caller, _ := UserFromContext(r.Context())
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", "Cuerpo inválido")
		return
	}
	tenantID, err := targetTenant(caller, req.TenantID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "tenant_required", "El super admin debe indicar tenantId")
		return
	}
	u, err := svc.CreateUser(r.Context(), tenantID, req.Email, req.Password, req.Name)
	switch {
	case errors.Is(err, ErrValidation):
		writeError(w, http.StatusBadRequest, "validation_error", "Email y contraseña (≥ 8) son requeridos")
	case errors.Is(err, ErrTenantNotFound):
		writeError(w, http.StatusNotFound, "tenant_not_found", "El negocio no existe")
	case errors.Is(err, ErrEmailTaken):
		writeError(w, http.StatusConflict, "email_taken", "El email ya está registrado")
	case err != nil:
		writeError(w, http.StatusInternalServerError, "internal", "No se pudo crear el usuario")
	default:
		writeJSON(w, http.StatusCreated, map[string]any{"user": u})
	}
}

func (svc *Service) handleListUsers(w http.ResponseWriter, r *http.Request) {
	caller, _ := UserFromContext(r.Context())
	tenantID, err := targetTenant(caller, r.URL.Query().Get("tenantId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "tenant_required", "El super admin debe indicar tenantId")
		return
	}
	users, err := svc.ListUsers(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "No se pudieron listar los usuarios")
		return
	}
	if users == nil {
		users = []User{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": users})
}

// targetTenant resuelve el negocio sobre el que opera la petición: el del usuario
// (si pertenece a uno) o el indicado explícitamente por el super admin global.
func targetTenant(caller User, explicit string) (string, error) {
	if caller.TenantID != nil {
		return *caller.TenantID, nil
	}
	if explicit == "" {
		return "", ErrTenantRequired
	}
	return explicit, nil
}
